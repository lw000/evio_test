syntax = "proto3";

package protocol; //包名

message RequestChat {
    string uid = 1;
    string msg = 2;
}

message ResponseChat {
    string msg = 2;
}