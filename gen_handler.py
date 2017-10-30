#!/usr/bin/env python


import subprocess
import shutil
import string
import shlex
import sys
import re
import os
import os.path

HANDLER_TPL = """

package handler

import (
	"gohive/internal/pb/msg"
)

{handlers}

"""

def gen_handlers(msgs, path):
    for p in msgs:
        cate = p[0]
        items = p[1]
        fullpath = os.path.join(path, cate) + ".go"
        ctx = ""
        for item in items:
            msg_name = item[0]
            msg_op = item[1]

            if msg_op == "Req":
                ctx += "func On%sReq(req *msg.%s_Req) *msg.%s_Res {\n" % (msg_name, msg_name, msg_name)
                ctx += "res := &msg.%s_Res{}\n" % (msg_name)
                ctx += "return res\n"
                ctx += "}\n\n"
            elif msg_op == "Res":
                ctx += "func On%sRes(res *msg.%s_Res) {\n" % (msg_name, msg_name)
                ctx += "}\n\n"
        #print ctx
        #print fullpath
        if not os.path.exists(fullpath):
            ctx = string.replace(HANDLER_TPL, "{handlers}", ctx)
            with open(fullpath, "w") as fp:
                fp.write(ctx)
                fp.close()

                # gofmt
                args = shlex.split("gofmt -w %s" % (fullpath))
                #subprocess.Popen(args).wait()
                subprocess.call(args)


CASE_TPL = """
case "msg.{name}.{op}":
    res = handler.On{name}{op}(m.(*msg.{name}_{op}))
    break
"""

DISPATCHER_TPL = """
package game

import (
	//log "github.com/Sirupsen/logrus"
	//"gohive/server/game/entity"
	"gohive/internal/pb/msg"
	"gohive/server/game/handler"
	"github.com/golang/protobuf/proto"
)

func Dispatch(api string, m proto.Message) (res proto.Message) {
	switch api {
        {case_ctx}
	default:
		break
	}
	return
}

"""

def _f(name):
    strs = []
    for s in name.split("_"):
        strs.append(s.capitalize())
    return "".join(strs)

def gen_dispatcher(msgs, dst_path):
    case_ctx = ""
    for p in msgs:
        items = p[1]
        for item in items:
            msg_name = item[0]
            msg_op = item[1]
            if msg_op == "Req":
                curr = string.replace(CASE_TPL, "{name}", msg_name)
                case_ctx += string.replace(curr, "{op}", msg_op)

    ctx = string.replace(DISPATCHER_TPL, "{case_ctx}", case_ctx)

    fullpath = os.path.join(dst_path, "dispatcher") + ".go"
    with open(fullpath, "w") as fp:
        fp.write(ctx)
        fp.close()

        # gofmt
        args = shlex.split("gofmt -w %s" % (fullpath))
        subprocess.call(args)

def get_msgs(src):
    proto_files = {}
    msg_types = []
    for path, dirs, files in os.walk(src, followlinks=True):
        for f in files:
            name, ext = os.path.splitext(f)
            if ext != ".proto":
                continue
            fullpath = os.path.join(path, f)

            protos = proto_files.get(path)
            if protos == None:
                protos = []
            protos.append(fullpath)
            proto_files[path] =  protos

            strs = path.split('/')
            if strs[len(strs) - 1] == "msg":
                msg_types.append((name,  parse_msg(fullpath)))
    return msg_types
    #gen_msg_dispatcher(msg_types, dispatcher_path)

def usage():
    print '%s [src_path] [dispatcher_path]' % (sys.argv[0])


MSG_BEGIN = r"\s*(message|enum)\s+(\w+)\s+\{"
MSG_END = r"\s*\}"

COMMENT_BEGIN = r"\s*\/\*"
COMMENT_END = r".*\*\/.*"
COMMENT = r"(\s*)\/\/.*"


def parse_msg(fullpath):
    stacks = []
    comment = False
    handlers = ""
    msgs = []
    with open(fullpath, "r") as fp:
        lines = fp.readlines()
        for l in lines:
            if re.compile(COMMENT).match(l) != None:
                continue

            if re.compile(COMMENT_BEGIN).match(l) != None:
                comment = True

            if re.compile(COMMENT_END).match(l) != None:
                comment = False
                continue

            if comment:
                continue

            r = re.compile(MSG_BEGIN)
            result = r.match(l)
            if result != None:
                stacks.append(l)

            r = re.compile(MSG_END)
            result = r.match(l)
            if result != None:
                s = stacks.pop()
                if len(stacks) > 0:
                    b = stacks[0]
                    res1 = re.compile(r"\s*(\w+)\s+(\w+)\s+").match(s)
                    res2 = re.compile(r"\s*(\w+)\s+(\w+)\s+").match(b)
                    op = res1.group(2)
                    name = res2.group(2)
                    if op != "Req" and op != "Res":
                        continue
                    msgs.append((name, op))
        fp.close()
    return msgs

if __name__ == "__main__":
    if len(sys.argv) != 4:
        usage()
        sys.exit(1)

    proto_source = sys.argv[1]
    dispatcher_path = sys.argv[2]
    handler_path = sys.argv[3]
    msgs = get_msgs(proto_source)
    gen_dispatcher(msgs, dispatcher_path)
    gen_handlers(msgs, handler_path)

