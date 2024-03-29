# Change these variables as necessary.
cmdName := kamcmd
cmdArg := cnxcc.active_clients
mysqlUser := root
mysqlPass := ANSKk08aPEDbFjDO
mysqlHost := 127.0.0.1
mysqlPort := 3306
mysqlDatabase := local


build:
	go build -o exec-cmd . 

exec:
	./exec-cmd \
		-cmdName=${cmdName} \
		-cmdArg=${cmdArg} \
		--mysqlUser=${mysqlUser} \
		--mysqlPass=${mysqlPass} \
		--mysqlHost=${mysqlHost} \
		--mysqlPort=${mysqlPort} \
		--mysqlDatabase=${mysqlDatabase}

run:
	go run . \
		-cmdName=${cmdName} \
		-cmdArg=${cmdArg} \
		--mysqlUser=${mysqlUser} \
		--mysqlPass=${mysqlPass} \
		--mysqlHost=${mysqlHost} \
		--mysqlPort=${mysqlPort} \
		--mysqlDatabase=${mysqlDatabase} \