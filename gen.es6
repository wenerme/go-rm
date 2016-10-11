function GenerateHashSetVar(name, ret, call) {
    let s = `
int HashSetVar(RedisModuleKey *key, int flags,int argc, intptr_t argv[]){
    switch(argc){
`

    for (var i = 0; i < 20; i++) {
        s += `case ${i}: return RedisModule_HashSet(key, flags`
        for (var j = 0; j < i; j++) {
            s += `,argv[${j}]`
        }
        s += `);\n`
    }

    s += `
    default:
        return REDISMODULE_ERR;
    }
}`
    return s
}
function GenerateHashGetVar(name, ret, call) {
    let s = `
int HashGetVar(RedisModuleKey *key, int flags,int argc, intptr_t argv[]){
    switch(argc){
`

    for (var i = 0; i < 20; i++) {
        s += `case ${i}: return RedisModule_HashGet(key, flags`
        for (var j = 0; j < i; j++) {
            s += `,argv[${j}]`
        }
        s += `);\n`
    }

    s += `
    default:
        return REDISMODULE_ERR;
    }
}`
    return s
}


// console.log(GenerateHashSetVar())
console.log(GenerateHashGetVar())
