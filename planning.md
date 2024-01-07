#### will use the mmap : 

## hirarky 

## Bucket / Collection
## Node   
## Inode


## Page
# corner case :
1. if the value of inod > pageSize then how to handle the split

```go 
type PgId unit32 
type PgType unit8

const (
   IPage PgType = iota // index page 
   DPage PgType  // data page
)

type Inode struct{
   key []byte
   Value PgId // this is dataPage  
}

type node struct{
   PgId PgId // this is IndexPage
   IsLeaf bool 
   Key []byte
   Parent PgId // this is IndexPage
   Childs PgId // this is IndexPage
   Inodes []Inode
}


type Page struct{
   PgType
   size uint32
   PgId PgId
}

type IndexPage struct{
   key []byte
}

type DataPage struct{
   key []byte
   value []byte
}
```