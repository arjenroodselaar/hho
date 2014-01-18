
if exists("b:current_syntax")
  finish
endif

syn keyword hhasOpcode Nop PopA PopC PopV PopR Dup Box Unbox BoxR BoxRNop UnboxR UnboxRNop Null True False
syn keyword hhasOpcode NullUninit Int Double String Array NewArray NewArrayReserve NewPackedArray NewStructArray AddElemC AddElemV AddNewElemC AddNewElemV NewCol ColAddElemC
syn keyword hhasOpcode ColAddNewElemC Cns CnsE CnsU ClsCns ClsCnsD File Dir Concat Abs Add Div Mod Sqrt Strlen
syn keyword hhasOpcode Xor Not Same NSame Eq Neq Lt Lte Gt Gte Shl Shr Floor Ceil CastBool
syn keyword hhasOpcode CastInt CastDouble CastString CastArray CastObject InstanceOf InstanceOfD Print Clone Exit Fatal Jmp JmpNS JmpZ JmpNZ
syn keyword hhasOpcode Switch SSwitch RetC RetV Unwind Throw CGetL CGetL2 CGetL3 PushL CGetN CGetG CGetS VGetL VGetN
syn keyword hhasOpcode VGetG VGetS AGetC AGetL IssetC IssetL IssetN IssetG IssetS EmptyL EmptyN EmptyG EmptyS IsTypeC IsTypeL
syn keyword hhasOpcode SetL SetN SetG SetS SetOpL SetOpN SetOpG SetOpS IncDecL IncDecN IncDecG IncDecS BindL BindN BindG
syn keyword hhasOpcode BindS UnsetL UnsetN UnsetG FPushFunc FPushFuncD FPushFuncU FPushObjMethod FPushObjMethodD FPushClsMethod FPushClsMethodF FPushClsMethodD FPushCtor FPushCtorD DecodeCufIter
syn keyword hhasOpcode FPushCufIter FPushCuf FPushCufF FPushCufSafe CufSafeArray CufSafeReturn FPassC FPassCW FPassCE FPassV FPassVNop FPassR FPassL FPassN FPassG
syn keyword hhasOpcode FPassS FCall FCallArray FCallBuiltin BaseC BaseR BaseL BaseLW BaseLD BaseLWD BaseNC BaseNL BaseNCW BaseNLW BaseNCD
syn keyword hhasOpcode BaseNLD BaseNCWD BaseNLWD BaseGC BaseGL BaseGCW BaseGLW BaseGCD BaseGLD BaseGCWD BaseGLWD BaseSC BaseSL BaseH ElemC
syn keyword hhasOpcode ElemL ElemCW ElemLW ElemCD ElemLD ElemCWD ElemLWD ElemCU ElemLU NewElem PropC PropL PropCW PropLW PropCD
syn keyword hhasOpcode PropLD PropCWD PropLWD PropCU PropLU CGetElemC CGetElemL VGetElemC VGetElemL IssetElemC IssetElemL EmptyElemC EmptyElemL SetElemC SetElemL
syn keyword hhasOpcode SetOpElemC SetOpElemL IncDecElemC IncDecElemL BindElemC BindElemL UnsetElemC UnsetElemL VGetNewElem SetNewElem SetOpNewElem IncDecNewElem BindNewElem CGetPropC CGetPropL
syn keyword hhasOpcode VGetPropC VGetPropL IssetPropC IssetPropL EmptyPropC EmptyPropL SetPropC SetPropL SetOpPropC SetOpPropL IncDecPropC IncDecPropL BindPropC BindPropL UnsetPropC
syn keyword hhasOpcode UnsetPropL CGetM VGetM FPassM IssetM EmptyM SetM SetWithRefLM SetWithRefRM SetOpM IncDecM BindM UnsetM IterInit IterInitK
syn keyword hhasOpcode WIterInit WIterInitK MIterInit MIterInitK IterNext IterNextK WIterNext WIterNextK MIterNext MIterNextK IterFree MIterFree CIterFree IterBreak Incl
syn keyword hhasOpcode InclOnce Req ReqOnce ReqDoc Eval DefFunc DefCls NopDefCls DefCns DefTypeAlias This BareThis CheckThis InitThisLoc StaticLoc
syn keyword hhasOpcode StaticLocInit Catch ClassExists InterfaceExists TraitExists VerifyParamType Self Parent LateBoundCls NativeImpl IncStat AKExists CreateCl Idx ArrayIdx
syn keyword hhasOpcode AssertTL AssertTStk AssertObjL AssertObjStk PredictTL PredictTStk BreakTraceHint CreateCont ContEnter ContSuspend ContSuspendK UnpackCont ContRetC ContCheck ContRaise
syn keyword hhasOpcode ContValid ContKey ContCurrent ContStopped ContHandle AsyncAwait AsyncESuspend AsyncWrapResult AsyncWrapException

syn match hhasFunction  "[a-zA-Z0-9_-]\+"
syn match hhasVariable  "\$[a-zA-Z0-9_-]\+"
syn match hhasDirective "^\.[a-zA-Z]\+" nextgroup=hhasFunction
syn match hhasComment "#.*$"
syn match hhasLabel   "\s\+[a-zA-Z0-9_]\+:"
syn match hhasString  "\".*\""
syn match hhasNumber  "\d\+"

let b:current_syntax = "hhas"

hi link hhasFunction  Function
hi link hhasOpcode    Statement
hi link hhasDirective Comment
hi link hhasComment   Comment
hi link hhasLabel     Operator
hi link hhasString    Constant
hi link hhasNumber    Constant
hi link hhasVariable  Type
