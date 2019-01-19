let argc = 21;

function GenerateHashSetVar(name, ret, call) {
    let s = `
int HashSetVar(RedisModuleKey *key, int flags,int argc, intptr_t argv[]){
    switch(argc){
`;

    for (var i = 0; i < argc; i++) {
        s += `case ${i}: return RedisModule_HashSet(key, flags`;
        for (var j = 0; j < i; j++) {
            s += `,argv[${j}]`
        }
        s += `);\n`
    }

    s += `
    default:
        return REDISMODULE_ERR;
    }
}`;
    return s
}

function GenerateHashGetVar(name, ret, call) {
    let s = `
int HashGetVar(RedisModuleKey *key, int flags,int argc, intptr_t argv[]){
    switch(argc){
`;

    for (var i = 0; i < argc; i++) {
        s += `case ${i}: return RedisModule_HashGet(key, flags`;
        for (var j = 0; j < i; j++) {
            s += `,argv[${j}]`
        }
        s += `);\n`
    }

    s += `
    default:
        return REDISMODULE_ERR;
    }
}`;
    return s
}


function GenerateCallVar(name, ret, call) {
    let s = `
RedisModuleCallReply * CallVar(RedisModuleCtx *key, const char *cmdname, const char *fmt, const int argc,const intptr_t argv[]){
    switch(argc){
`;

    for (var i = 0; i < argc; i++) {
        s += `case ${i}: return RedisModule_Call(key, cmdname, fmt`;
        for (var j = 0; j < i; j++) {
            s += `,argv[${j}]`
        }
        s += `);\n`
    }

    s += `
    default:
        errno=EINVAL;
        return NULL;
    }
}`;
    return s;
}

const all = [
    GenerateCallVar(),
    '/* Hash */',
    GenerateHashSetVar(),
    GenerateHashGetVar(),
];

// console.log(GenerateHashSetVar())
// console.log(GenerateHashGetVar())
// console.log(GenerateCallVar())
const fs = require('fs');
let f = fs.readFileSync('rm/varargs.h').toString();

// region Generated

f = f.replace(/(?<=region Generated)(.*)(?=\n[^\n]+endregion)/s, all.join('\n'));


fs.writeFileSync('rm/varargs.h', f);
