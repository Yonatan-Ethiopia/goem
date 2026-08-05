package assert

import (
        "testing"
        "strings"
    )
    
func Equal[T comparable](t *testing.T, actual, expected T){
    t.Helper()
    
    if actual != expected{
        t.Errorf("got: %d want %d", actual, expected)
    }
}

func StringContains(t *testing.T, actual, expectedString string){
    t.Helper()
    
    if !strings.Contains(actual, expectedString){
        t.Errorf("got: %q expected to contain: %q", actual, expectedString)
    }
}
