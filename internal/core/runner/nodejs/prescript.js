const argv = process.argv

const koffi = require('koffi')
const lib = koffi.load('./var/sandbox/sandbox-nodejs/nodejs.so')
const kozmoSeccomp = lib.func('void KozmoSeccomp(int, int, bool)')

const uid = parseInt(argv[2])
const gid = parseInt(argv[3])

const options = JSON.parse(argv[4])

kozmoSeccomp(uid, gid, options['enable_network'])

