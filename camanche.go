package camanche

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"github.com/TateHanawalt/goutil"
	"github.com/TateHanawalt/camanche/opts"
)

type CallFn func(ArgType) error

const cmdlbl = "CMD"
const urilbl = "URI"
const flaglbl = "FLAG"
const paramlbl = "PARAM"

type KV struct {
	Key string `yaml:"key"`
	Val string `yaml:"val"`
}
type KVSet []KV

// name + location
type Uri struct {
	Name string
	Location string
}

type ArgType struct {
	ArgCount   int
	curindex   int
	Cmd        string
	Exec       string
	Flags      []string
	lastparams map[string]string
	Ledger     []KV
	Params     map[string][]string
	URI				 []Uri
	LastURI    *KV
}
func (a *ArgType) Shift() *ArgType {
	if a.curindex >= len(a.Ledger) {
		a.Cmd = ""
		return a
	}
	a.curindex = a.curindex + 1
	for a.curindex < len(a.Ledger) {
		kv := a.Ledger[a.curindex]
		switch kv.Key {
		case cmdlbl:
			a.Cmd = kv.Val
			return a
			break
		case paramlbl:
			indKey := strings.Split(kv.Val, ":")
			index, err := strconv.Atoi(indKey[0])
			if err != nil {
				panic(err)
			}
			if index < 0 {
				panic(fmt.Errorf("Camanche map index was lt 0..."))
			}
			ParamVal := a.Params[indKey[1]][index]
			a.lastparams[indKey[1]] = ParamVal
		case flaglbl:
			// flags
		case urilbl:
			a.LastURI = &kv
		default:
			fmt.Printf("\n Non Cmd or Param: %s\n", kv.Key)
		}
		a.curindex += 1
	}
	if a.curindex >= len(a.Ledger) {
		a.Cmd = ""
	}
	return a
}
func (a *ArgType) ReadREm() {
	if a.curindex >= len(a.Ledger) {
		a.Cmd = ""
		return
	}
	a.curindex = a.curindex + 1
	for a.curindex < len(a.Ledger) {
		kv := a.Ledger[a.curindex]
		switch kv.Key {
		case paramlbl:
			indKey := strings.Split(kv.Val, ":")
			index, err := strconv.Atoi(indKey[0])
			if err != nil {
				panic(err)
			}
			if index < 0 {
				panic(fmt.Errorf("Camanche map index was lt 0..."))
			}
			ParamVal := a.Params[indKey[1]][index]
			a.lastparams[indKey[1]] = ParamVal
		case flaglbl:
			// flags
		default:
			fmt.Printf("\n Non Cmd or Param: %s\n", kv.Key)
		}
		a.curindex += 1
	}
}
func (a *ArgType) RemCmds() int {
	rem := 0
	iter := a.curindex + 1
	for iter < len(a.Ledger) {
		kv := a.Ledger[iter]
		if kv.Key == cmdlbl {
			rem += 1
		}
		iter += 1
	}
	return rem
}
func (a ArgType) 	NumCmds() int {
	count := 0
	for i := range a.Ledger {
		if a.Ledger[i].Key == cmdlbl {
			count += 1
		}
	}
	return count
}
func (a ArgType) Next(oftype string) *KV {
	// tmpindex := a.curindex
	cpindex := a.curindex + 1
	for cpindex < len(a.Ledger) {
		kv := a.Ledger[cpindex]
		if kv.Key == oftype {
			return &kv
		}
	}
	return nil
}
func (a ArgType) NextURI() (string, error) {
	val := a.Next(urilbl)
	if val != nil {
		return val.Val, nil
	}
	return "", nil
}
func (a ArgType) lastnoftype(oftype string, n int) ([]KV, error) {
	var rset []KV
	index := a.curindex
	if index >= len(a.Ledger) {
		index = len(a.Ledger) - 1
	}
	for index >= 0 {
		if a.Ledger[index].Key == oftype {
			if len(rset) < n || n < 0 {
				// Add entry to front of array
				rset = append([]KV{a.Ledger[index]}, rset...)
			}
		}
		index -= 1
	}
	return rset, nil
}
func (a ArgType) LastNUri(cnt int) ([]string, error) {
	kvset, err := a.lastnoftype(urilbl, cnt)
	if err != nil {
		return nil, err
	}
	rset := []string{}
	for i := range kvset {
		rset = append(rset, kvset[i].Val)
	}
	return rset, nil
}
func (a *ArgType) ParamVal(ParamKey string) (string, bool) {
	if val, ok := a.lastparams[ParamKey]; ok {
		return val, true
	}
	return "", false
}
func (a *ArgType) TrueFlag(flag string) bool {
	if len(flag) < 1 {
		return false
	}
	for f := range a.Flags {
		if a.Flags[f] == flag {
			return true
		}
	}
	return false
}
func (a *ArgType) Print() {

	fmt.Println()
	fmt.Println(" Printing Camanche Arg Type: ")
	fmt.Println()
	fmt.Printf(" curindex  %d\n", a.curindex)
	fmt.Printf(" Cmd       %s\n", a.Cmd)
	fmt.Printf(" ArgCount  %d\n", a.ArgCount)
	fmt.Printf(" Flags Num %d\n", len(a.Flags))

	if len(a.Flags) > 0 {
		fmt.Println()
	}
	for index, i := range a.Flags {
		fmt.Printf("  flag %d - %v\n", index, i)
	}
	fmt.Println()

	fmt.Printf(" Params Num %d\n", len(a.Params))
	if len(a.Params) > 0 {
		fmt.Println()
	}
	for index, i := range a.Params {
		fmt.Printf("  key: %v,  values: %v\n", index, i)
	}
	fmt.Println()
	fmt.Printf(" Exec:       %v\n", a.Exec)
	fmt.Printf(" lastparams: %v\n", a.lastparams)

	fmt.Printf(" Ledger: ")
	fmt.Println()
	for _, k := range a.Ledger {
		fmt.Printf("         %v - %v\n", goutil.MakeW(k.Key, 6), k.Val)
	}
	fmt.Println()
}

type Opts struct {
	ArgType
}
func (o Opts) Verbose() bool {
	return o.TrueFlag("v")
}
func (o Opts) All() bool {
	return o.TrueFlag("a")
}
func (o Opts) Force() bool {
	return o.TrueFlag("force")
}
func (o Opts) ShowPrivate() bool {
	return o.TrueFlag("show")
}
func (o Opts) Json() bool {
	return o.TrueFlag("j")
}
func (o Opts) Yaml() bool {
	return o.TrueFlag("y")
}
func (o Opts) Kvset() bool {
	return o.TrueFlag("kvs")
}
func (o Opts) Indent() int {
	if id, ok := o.Params["indent"]; ok && len(id) > 0 {
		idval := id[0]
		i, err := strconv.Atoi(idval)
		if err != nil {
			return 0
		}
		if i > 0 {
			return i
		}
	}
	return 0
}
func (o Opts) Help() bool {
	return o.TrueFlag("h")
}
func (o Opts) Flag(flagid string) bool {
	return o.TrueFlag(flagid)
}

type CMD struct {
	id string
	fn CallFn
	cmds map[string]*CMD
}
func (c CMD)  Id() string {
	return c.id
}
func (c *CMD) SetId(id string) {
	c.id = id
}
func (c *CMD) Add(newcmd *CMD) error {
	if newcmd == nil {
		return fmt.Errorf("Parameter newcmd is nil...")
	}
	command := newcmd.Id()
	if len(command) < 1 {
		return fmt.Errorf("Command length must be gt 0")
	}
	if _, found := c.cmds[command]; found {
		return fmt.Errorf("Command %s already specified for %s", command, c.id)
	}
	if c.cmds == nil {
		c.cmds = make(map[string]*CMD)
	}
	c.cmds[command] = newcmd
	return nil
}
func (c *CMD) Parse(args ArgType) error {
	if c == nil {
		return fmt.Errorf("Camanche CMD pointer receiver is nil in Camanche.Parse(args ArgType) error")
	}

	if c.id != "root" {
		args.Shift()
	}
	command := args.Cmd
	if len(command) < 1 || len(c.cmds) < 1 {
		if c.fn != nil {
			return c.fn(args)
		}
		return fmt.Errorf("\n Args exausted at non handler node %s", c.id)
	}
	if next, found := c.cmds[command]; found {
		return next.Parse(args)
	}
	return fmt.Errorf("No handlers found in %s for %s", c.id, command)
}
func (c *CMD) MkAdd(id string, fn CallFn) (*CMD, error) {
	addcmd, err := NewCMD(id, fn)
	if err != nil {
		return nil, err
	}
	if addcmd == nil {
		return nil, fmt.Errorf("NewCMD returnd nil...")
	}
	err = c.Add(addcmd)
	return addcmd, err
}
func (c *CMD) MkAddNR(id string, fn CallFn) (error) {
	mk, err := c.MkAdd(id, fn)
	if err != nil {
		return err
	}
	if mk == nil {
		return fmt.Errorf("returned cmd node was nil in MkAddNR")
	}
	return nil
}





func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1
}
func NewCMD(id string, self CallFn) (*CMD, error) {
	root := &CMD{}
	root.SetId(id)
	root.fn = self
	return root, nil
}
func ToFix(args []string) []string {
	cpyset := make([]string, len(args))
	copy(cpyset, args)
	fixedSet := make([]string, len(args))
	loadImmediate := false
	for len(cpyset) > 0 {
		arg0 := cpyset[0]
		cpyset = cpyset[1:]
		arg0 = strings.TrimSpace(arg0)
		arg0 = strings.ReplaceAll(arg0, "\"", "")
		if loadImmediate {
			loadImmediate = false
			fixedSet = append(fixedSet, arg0)
			continue
		}
		split := strings.Split(arg0, "--")
		if len(split) == 1 {
			if len(split[0]) < 1 {
				continue
			}
			if split[0][0] == '-' {
				fixedSet = append(fixedSet, "-")
				pushval := split[0][1:]
				fixedSet = append(fixedSet, string(pushval))
				continue
			}
			appendVal := split[0]
			fixedSet = append(fixedSet, appendVal)
			continue
		}
		if len(split) != 2 {
			panic(fmt.Errorf(" split by '--' set length not 1 or 2 - this is an error"))
		}
		fixedSet = append(fixedSet, "--")
		kvset := strings.Split(split[1], "=")
		fixedSet = append(fixedSet, kvset[0])
		if len(kvset) == 2 {
			fixedSet = append(fixedSet, kvset[1])
			continue
		}
		if len(cpyset) < 1 {
			panic(fmt.Errorf("Args reached end of control flow expecting a value not nothing"))
		}
		loadImmediate = true
	}
	return fixedSet
}
func Parse() ArgType {
	rettype := ArgType{}
	rettype.Params     = make(map[string][]string)
	rettype.lastparams = make(map[string]string)

	args := os.Args
	rettype.Exec = args[0]

	args = args[1:]
	fixedSet := ToFix(args)

	cpyset := make([]string, len(fixedSet))
	copy(cpyset, fixedSet)
	inflag := false
	inkey := false
	invalue := false
	setKey := ""

	for len(cpyset) > 0 {
		arg0 := cpyset[0]
		cpyset = cpyset[1:]
		arg0 = strings.TrimSpace(arg0)

		if len(arg0) < 1 {
			continue
		}

		if inflag {
			inflag = false
			rettype.Flags = append(rettype.Flags, arg0)
			node := KV{}
			node.Key = flaglbl
			node.Val = arg0
			rettype.Ledger = append(rettype.Ledger, node)
			continue
		}

		if inkey {
			inkey = false
			invalue = true
			setKey = arg0
			continue
		}

		if invalue {
			invalue = false
			Value := arg0
			node := KV{}
			node.Key = paramlbl
			node.Val = setKey
			index := -1
			if _, ok := rettype.Params[setKey]; ok {
				valAt := indexOf(setKey, rettype.Params[setKey])
				if valAt < 0 {
					// Value not found in map, append, get index and record
					rettype.Params[setKey] = append(rettype.Params[setKey], Value)
					index = len(rettype.Params[setKey]) - 1
				} else {
					index = valAt
				}
			} else {
				initset := []string{}
				initset = append(initset, Value)
				rettype.Params[setKey] = initset
				index = 0
			}
			node.Val = fmt.Sprintf("%d:%s", index, node.Val)
			rettype.Ledger = append(rettype.Ledger, node)
			continue
		}

		if arg0 == "-" {
			inflag = true
			continue
		}

		if arg0 == "--" {
			inkey = true
			continue
		}

		node := KV{}
		node.Key = cmdlbl

		if strings.Contains(arg0, ":") {
			node.Key = urilbl
		}

		node.Val = arg0
		rettype.Ledger = append(rettype.Ledger, node)
		cmdcnt := rettype.NumCmds()

		if node.Key == cmdlbl {
			if cmdcnt == 1 {
				rettype.Cmd = arg0
			}
		}
	}
	rettype.ArgCount = len(rettype.Ledger)
	return rettype
}
func RemCmds() int {
	return 0
}
func OptsFromArg(arg *ArgType) opts.CommonOpts {
	return Opts{*arg}
}
