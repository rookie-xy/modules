package collector

import "github.com/rookie-xy/modules/agents/log/match"

// MatchAny checks if the text matches any of the regular expressions
func MatchAny(matchers []match.Matcher, text string) bool {
    for _, m := range matchers {
        if m.MatchString(text) {
            return true
        }
    }

    return false
}
