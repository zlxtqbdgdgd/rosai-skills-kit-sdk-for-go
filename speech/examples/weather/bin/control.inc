#!/bin/bash

#file:includes.sh

VERSION="1.0.3.0"

#Success
AOS_OK=0

#faile
AOS_ERR_FAILED_OPERATION=1

#already started
AOS_ERR_STARTED=2

#not started
AOS_ERR_NOT_STARTED=3

#timeout
AOS_ERR_TIMEOUT=4

#file not found
AOS_ERR_FILE_NOTFOUND=5

#err access
AOS_ERR_ACCESS=6

#err para
AOS_ERR_PARAMATER=7

#log not rotate,but module is ok
AOS_ERR_NOT_WORKING=8


function file_exist
{
    local file_name
    if [ $# -ne 1 ];then
        return 1
    fi
    file_name=$1
    if [ -e $file_name -a ! -d $file_name ];then
        return 0
    else
        return 1
    fi
}
function file_executable
{
    local file_name
    if [ $# -ne 1 ];then
        return 1
    fi
    file_name=$1
    if [ -x $file_name -a ! -d $file_name ];then
        return 0
    else
        return 1
    fi
}
function check_proc_exist
{
    local proc_name proc_port port_name
    if [ $# -eq 2 -a $2 -gt 0 2>/dev/null ];then
        proc_name=$1
        proc_port=$2
        port_name=`/usr/sbin/lsof -P -i :$proc_port |grep "*:$proc_port\>"|awk '{print $1}'|uniq`
        if [ "$proc_name" = "$port_name" ];then
            return 0
        else
            return 1
        fi
        elif [ $# -eq 1 ];then
        proc_name=$1
        pstree -pl|grep "\-${proc_name}([0-9].*)" &>/dev/null
        if [ $? -eq 0 ];then
            return 0
        else
            return 1
        fi
    else
        return 1
    fi
}
function check_proc_num
{
    local proc_num proc_name real_num
    local default_num=3
    if [ $# -eq 2 -a $2 -gt 0 2>/dev/null ];then
        proc_name=$1
        proc_num=$2
    elif [ $# -eq 1 ];then
        proc_name=$1
        proc_num=$default_num
    else
        return 1
    fi
    real_num=`pstree -pl|grep "\-[{]\{0,1\}${proc_name}[}]\{0,1\}([0-9].*)"|wc -l |awk '{print$1}'`
    #real_num=` -pl|grep "\-$proc_name([0-9].*)"|wc -l|awk '{print $1}'`
    if [ $real_num -ge $proc_num ];then
        return 0
    else
        return 1
    fi
}

function check_log_roll
{
    local time_limit log_name i this_time last_time
    if [ $# -ne 2 ];then
        return 1
    fi
    if [ ! -r $1 ];then
        return 1
    fi
    if [ $2 -gt 0 2>/dev/null ];then
        time_limit=$2
    else
        return 1
    fi
    log_name=$1
    for((i=0,this_time=`date +%s -r $log_name`,last_time=this_time;i<=time_limit;i++))
    do
        if [ $this_time -gt $last_time ];then
            return 0
        fi
        sleep 1
        last_time=$this_time
        this_time=`date +%s -r $log_name`
    done
    return 1
}


#example: kill_proc "as_r"
#         will kill supervise & as_r
#propose: kill process of "as_r" and all parents of its
#
function kill_proc
{
    local proc_name all_id proc_id father_id i
    if [ $# -ne 2 ];then
        return 1
    fi
    for arg in $@
    do
        proc_name=$arg

            father_id=`ps -eo ppid,comm|grep " $proc_name\>"|awk '{print $1}'|head -n 1|grep -v "\<1\>"`

            if [ ! -z "$father_id" ];then
                kill -9 $father_id 2>/dev/null
            fi

        all_id=`ps -eo ppid,pid,comm|grep " $proc_name\>"|awk '{printf $1"\n"$2"\n"}'|sort -n|uniq|grep -v "\<1\>"`
        if [ -z "$all_id" ];then
    #return 1
    continue
        fi

        for proc_id in $all_id
        do
            kill $proc_id 2>/dev/null
            sleep 0.1
            kill -9 $proc_id 2>/dev/null
        done
    done
    #sleep 1
    #killall -9 $proc_name 2>/dev/null
    #sleep 1
    #killall -9 $proc_name 2>/dev/null
    return 0
}


