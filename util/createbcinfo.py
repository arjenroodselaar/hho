#!/usr/bin/env python2.7

import re

# Nop []  ->  []
# PopA [A]  ->  []
# PopC [C]  ->  []
# PopV [V]  ->  []
# PopR [R]  ->  []
# Dup [C:<T>]  ->  [C:<T> C:<T>]

regex = re.compile("^(.*)\s+\[(.*?)\]\s+->\s+\[(.*?)\]$")

header = """
func LookupStackDelta(bc string) int {
    switch bc {"""

footer = """    default:
        panic("Invalid bytecode passed to LookupStackDelta")
    }
    return 0
}"""

def print_go_switcher(op, stack_mod):
    print "    case \"{}\":\n        return {}".format(op, stack_mod)

print header
for i in open('bytecodes','r').readlines():
    op, lstack, rstack = regex.match(i).groups()
    stack_before = len(lstack.split())
    stack_after = len(rstack.split())
    stack_mod = stack_after - stack_before
    print_go_switcher(op, stack_mod)

print footer
