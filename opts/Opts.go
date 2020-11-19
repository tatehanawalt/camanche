package opts

import (
  "fmt"
  "strconv"
)

type boolfrom func(string) bool
type boolflagprovider interface {
  TrueFlag(string) bool
}
type strparamprovider interface {
  ParamVal(string) string
}

type verbosedecider interface {
  Verbose() bool
}
type jsondecider interface {
  Json() bool
}
type CommonOpts interface {
  verbosedecider
  jsondecider

  Force() bool
	All() bool
	ShowPrivate() bool
	Yaml() bool
	Kvset() bool
	Help() bool
	Indent() int
	Flag(string) bool
}

func CheckFlags(bfg boolflagprovider, fset []string) bool {
  for _, flag := range fset {
    if bfg.TrueFlag(flag) {
      return true
    }
  }
  return false
}
func Verbose(target interface{}) bool {
  res := false
  if vbif, isvp := target.(verbosedecider); isvp {
    res = vbif.Verbose()
  } else if bfg, ifp := target.(boolflagprovider); ifp {
    res = CheckFlags(bfg, []string{"v", "V", "verbose", "Verbose"})
  }
  return res
}
func BoolFlag(target interface{}) boolfrom {
  if bfg, ifp := target.(boolflagprovider); ifp {
    return bfg.TrueFlag
  }
  return func(val string) bool {
    return false
  }
}
func Json(target interface{}) bool {
  if jsif, isvp := target.(jsondecider); isvp {
    return jsif.Json()
  } else if bfg, ifp := target.(boolflagprovider); ifp {
    return CheckFlags(bfg, []string{"j", "J", "json", "Json"})
  } else {
    fmt.Printf(" Got non types for %v\n", target)
  }
  return false
}
func Yaml(target interface{}) bool {
  if bfg, ifp := target.(boolflagprovider); ifp {
    return CheckFlags(bfg, []string{"y", "Y", "yml", "Yml"})
  } else {
    fmt.Printf(" Got non types for %v\n", target)
  }
  return false
}
func Indent(target interface{}) int {
  if spp, isp := target.(strparamprovider); isp {
    indentstr := spp.ParamVal("indent")
    i, err := strconv.Atoi(indentstr)
    if err != nil {
      fmt.Printf("\n Opts err: %v\n", err)
      return 0
    }
    if i < 0 {
      return 0
    }
    return i;
  }
  return 0
}
