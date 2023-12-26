#!/bin/sh

go mod download

i=0
declare -A array=(
	["healthcheck"]="healthcheck"
  ["memo-post"]="memo/post"
  ["memo-get"]="memo/get"

)

i=0
for key in ${!array[@]}; do
	echo "${array[$key]}をコンパイル開始"
	GOOS=linux GOARCH=arm64 go build -tags lambda.norpc -o bin/${array[$key]}/bootstrap src/${array[$key]}/main.go
	if [ $? -ne 0 ]; then
		echo "${array[$key]}のコンパイル中にエラーが発生"
		exit 1
	fi
	echo "${array[$key]}をコンパイル終了"
	# zipファイルにbootstrapとyamlを格納
	echo "${array[$key]}をzipファイルに格納"
	zip -j bin/${key}.zip bin/${array[$key]}/bootstrap
	if [ $? -ne 0 ]; then
		echo "${array[$key]}のzipファイルに格納中にエラーが発生"
		exit 1
	fi
	echo "${array[$key]}をzipファイルに格納完了"
	let i++
done

