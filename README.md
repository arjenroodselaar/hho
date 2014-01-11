hho
===

A Golang to HHAS compiler

Running HHVM as 'as'
====================
./hhvm -v Eval.AllowHhas=true <file>

Plan of Attack
==============
Using ssadump as an example, we're going to do the same thing.
But instead of spitting out SSA, we're pooping out HHAS.

Bitchin'
