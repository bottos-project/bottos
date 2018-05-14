rem generate common/types
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\  ..\..\..\..\github.com\bottos-project\core\common\types\transaction.proto ..\..\..\..\github.com\bottos-project\core\common\types\basic-transaction.proto ..\..\..\..\github.com\bottos-project\core\common\types\block.proto

rem generate api
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\ --micro_out=:..\..\..\..\  ..\..\..\..\github.com\bottos-project\core\api\transaction.proto ..\..\..\..\github.com\bottos-project\core\api\basic-transaction.proto  ..\..\..\..\github.com\bottos-project\core\api\core-api.proto

rem generate tool/example
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=:..\..\..\..\ --micro_out=:..\..\..\..\   ..\..\..\..\github.com\bottos-project\core\tool\example\example.proto
