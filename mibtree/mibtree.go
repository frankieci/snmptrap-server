package mibtree

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	ID string
	Name string
	Desc string
	Children map[string]*Node
}

func NewMibTree() (*Node) {
	root := &Node {
		ID: ".",
		Name : "root",
		Desc : "root Node",
		Children: make(map[string]*Node,0),
	}
	return root
}

func (n *Node)AddNode(oid string, name string, desc string) (error) {
	if strings.HasPrefix(oid, ".") {
		oid = oid[1:]
	}
	oids := strings.Split(oid, ".")

	curNode := n
	for _, id := range oids {
		if _, err := strconv.ParseInt(id, 10, 64); err != nil {
			return err
		}

		if nextNode, ok := curNode.Children[id]; ok {
			curNode = nextNode
			continue
		} else {
			newNode := &Node {
				ID: id,
				Name : "",
				Desc : "",
				Children: make(map[string]*Node,0),
			}
			curNode.Children[id] = newNode
			curNode = newNode
			continue
		}
	}
	curNode.Name = name
	curNode.Desc = desc
	return nil
}


func (n *Node)FindNodeName(oid string) (name string, err error) {
	if strings.HasPrefix(oid, ".") {
		oid = oid[1:]
	}
	oids := strings.Split(oid, ".")
	curNode := n
	for i := 0; i < len(oids); i++ {
		id := oids[i];
		if _, err := strconv.ParseInt(id, 10, 64); err != nil {
			return "", errors.New(fmt.Sprintf("[%s]该id不是整型数字;%s",id, err))
		}

		if nextNode, ok := curNode.Children[id]; ok {
			if i == len(oids) - 1 {
				return nextNode.Name, nil
			} else {
				curNode = nextNode
				continue
			}
		} else {
			return curNode.Name + "." + strings.Join(oids[i:], "."), nil
		}

	}
	return "", errors.New(fmt.Sprintf("have not such oid[%s]", oid))
}

func (n *Node)Print(depth int64) {
	sep := "-"
	for  i := int64(0); i < depth; i++ {
		sep = sep + "-"
	}
	namepath := fmt.Sprintf("|%s[%s:%s,%s]", sep, n.ID,n.Name,n.Desc)
	for  j := int64(0); j < depth; j++ {
		namepath = "   " + namepath
	}
	fmt.Println(namepath)
	childrens := n.Children
	for _,v := range childrens {
		v.Print( depth + 1 )
	}
}


func (n *Node)LoadFile(filepath string) ( error) {
	fi, err := os.Open(filepath)
	if err != nil {
		return errors.New(fmt.Sprintf("开发文件报错：%s", err))
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	name := ""
	oid := ""
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		strAry := strings.Fields(string(a))
	    if len(strAry) != 2 {
	    	continue
		}else {
			name = strings.Trim(strAry[0],"\"")
			oid = strings.Trim(strAry[1],"\"")
		}
		n.AddNode(oid, name, "")
	}
	return nil
}