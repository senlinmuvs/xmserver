kill `ps -aux|grep xmserver | grep -v "grep --color" | head -n 1 | awk '{print $2}'`
