// v1.0
package httprouter

import "net/http"

var router = map[string]*Router {
	http.MethodGet: &Router{},
	http.MethodPost: &Router{},
	http.MethodDelete: &Router{},
	http.MethodConnect: &Router{},
	http.MethodHead: &Router{},
	http.MethodOptions: &Router{},
	http.MethodPatch: &Router{},
	http.MethodPut: &Router{},
	http.MethodTrace: &Router{},
	"ALL": &Router{},
}

type Params []Param

type Param struct {
	Key   string
	Value string
}

type Router struct {
	str string         		// 前缀字符串
	param Param				// 每个节点只能存一个
	childNode []*Router     // 该节点的所有子节点
}

type Handle func(w http.ResponseWriter, r *http.Request)

func New() {

}

func (r *Router)Get(path string, handle Handle) *Router {
	r.Handle(http.MethodGet ,path, handle)
	return r
}

func (r *Router)Handle(method, path string, handle Handle) {
	r.handle(method, path, handle)
}

func (Router)handle(method, path string, handle Handle) {
	//router[method]
}

//func (t *Router) Insert(str string) {
//	if len(t.childNode) == 0 {
//		t.childNode = append(t.childNode, &Router{str:str})
//	} else {
//		i := 0
//		for ;i < len(t.childNode);i++ {
//			// 可以挂载到这个节点的 子节点为 i 的节点上
//			if str[0] == t.childNode[i].str[0] {
//				// 寻找公共最长前缀
//				j := 0
//				for ; j < len(str) && j < len(t.childNode[i].str) && str[j] == t.childNode[i].str[j]; j++ {
//				}
//				if j < len(str) && j < len(t.childNode[i].str) {
//					tempN := &node{str:str[:j]}
//
//					t.childNode[i].str = t.childNode[i].str[:j]
//					tempN.childNode = append(tempN.childNode, t.childNode[i])
//					tempN.childNode = append(tempN.childNode, &node{str:str[j:]})
//
//					t.childNode[i] = tempN
//				} else if len(str) == len(t.childNode[i].str) && len(str) == j {
//					t.childNode[i].val = &str
//				} else if len(str) > len(t.childNode[i].str) {
//					t.childNode[i].Insert(str[j:])
//				} else {
//					tempN := &node{str:str}
//					tempN.childNode = append(tempN.childNode, t.childNode[i])
//					t.childNode[i] = tempN
//				}
//				return
//			}
//		}
//		// 挂载到这个节点上
//		t.childNode = append(t.childNode, &node{str:str})
//	}
//}
