
# validate
    import "github.com/zer0131/toolbox/base/validate"


# Demo
```
var params TestStruct
br := validate.BindFromJson(&params, jsonStr)
if !br.OK || len(br.FieldErrors) != 0 {
    return Reply(1, "参数出错", br.FieldErrors)
}
```
