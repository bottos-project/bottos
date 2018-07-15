rem generate common/types
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\  ..\..\..\..\github.com\bottos-project\bottos\common\types\transaction.proto ..\..\..\..\github.com\bottos-project\bottos\common\types\basic-transaction.proto ..\..\..\..\github.com\bottos-project\bottos\common\types\block.proto

rem generate api
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\ --micro_out=:..\..\..\..\  ..\..\..\..\github.com\bottos-project\bottos\api\transaction.proto ..\..\..\..\github.com\bottos-project\bottos\api\basic-transaction.proto  ..\..\..\..\github.com\bottos-project\bottos\api\chain.proto

rem generate tool/example
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\ --micro_out=:..\..\..\..\   ..\..\..\..\github.com\bottos-project\bottos\tool\example\example.proto
