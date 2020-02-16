package reformat

// import "github.com/stretchr/testify/assert"
// import "testing"

// type (
// 	atomWrapper struct {
// 		before string
// 		after  string
// 	}
// 	reformatConfig struct {
// 		line      int
// 		column    int
// 		indent    int
// 		spacing   int
// 		separator string
// 		wrapper   atomWrapper
// 		rule      interface{}
// 	}
// 	ruleMaxAtomsPerLine struct {
// 		max int
// 	}
// 	ruleFitLineLimit struct {
// 		lineLimit int
// 	}
// )

// func TestTreeReformat(t *testing.T) {
// 	before := `
// {"a":"b","c":{"d":"e","f":{"g":"h"}}}
// `
// 	cfg := reformatConfig{
// 		line:      1,
// 		column:    2,
// 		indent:    2,
// 		spacing:   1,
// 		separator: ",",
// 		wrapper:   atomWrapper{"{", "}"},
// 		rule:      ruleMaxAtomsPerLine{1},
// 	}
// 	after := `
// {
//   "a":"b",
//   "c":{
//     "d":"e",
//     "f":{
//       "g":"h",
//       "a":"x"
//     }
//   }
// }
// `
// 	cfg.rule = ruleMaxAtomsPerLine{2}
// 	after = `
// {
//   "a":"b",
//   "c": {
//     "d": "e",
//     "f": {
//       "g":"h", "a":"x"
//     }
//   }
// }
// `
// 	cfg.rule = ruleMaxAtomsPerLine{3}
// 	after = `
// {
//   "a":"b",
//   "c": {
//     "d": "e", "f":{"g":"h", "a":"x"}
//   }
// }
// `
// }
