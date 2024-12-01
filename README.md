# G-aB

## 代码规范
+ API层(api参数)
+ Model层(gorm模型)
+ Service层(数据库读取操作层)

## 层级传递
+ 每个模块分别在文件内建立一个struct,然后var _ = new (xxx)校验一次
最外层拿到这个struct,包裹在一个总的struct内,并用var xxxxxApp = new(xxxxStruct)定义一个值,后续别的模块使用就直接使用这个值内的值来使用方法
