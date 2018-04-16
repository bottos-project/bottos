rem generate common/types
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=plugins=micro:..\..\..\..\   ..\..\..\..\github.com\bottos-project\core\common\types\transaction.proto ..\..\..\..\github.com\bottos-project\core\common\types\block.proto

rem generate tool/example
protoc\bin\protoc.exe -I ..\..\..\..\ --go_out=plugins=micro:..\..\..\..\   ..\..\..\..\github.com\bottos-project\core\tool\example\example.proto
