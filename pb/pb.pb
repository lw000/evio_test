syntax = "proto3";

package protocol; //包名

message RequestChat {
    string uid = 1;
    string msg = 2;
    int64  t = 3;
}

message ResponseChat {
    string msg = 1;
    int64 t = 2;
}
