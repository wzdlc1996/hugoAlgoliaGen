package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_tomlParser(x *testing.T) {
	testStr := `
	key1                = false
	key2               = "zh_Hans"
	key3               = 5
	[permalnks]
		blog		= "asdfdsad"
	[aaa]
		test	= "a"
		test2 = "b"
	`
	refMap := map[string]string{
		"key1": "false",
		"key2": "zh_Hans",
		"key3": "5",
	}
	resMap := tomlParser(testStr)
	fmt.Println(resMap)
	for key, val := range refMap {
		if resMap[key] != val {
			x.Errorf("Failed at %s, expected %s, but get %s", key, val, resMap[key])
		}
	}

}

func Test_yamlParser(x *testing.T) {
	testStr := `
	key1                : false
	key2               : "zh_Hans"
	key3               : 5
	aaa :
		- pikachu
		- no!
	vvv
		test	: "a"
		test2 : "b"
	`
	refMap := map[string]string{
		"key1": "false",
		"key2": "zh_Hans",
		"key3": "5",
	}
	resMap := yamlParser(testStr)
	fmt.Println(resMap)
	for key, val := range refMap {
		if resMap[key] != val {
			x.Errorf("Failed at %s, expected %s, but get %s", key, val, resMap[key])
		}
	}

}

func TestConfigParser(t *testing.T) {
	input := "aaa"
	res, err := ConfigParser(input)
	if err == nil {
		t.Errorf("Failed at illegal input, expected non nil, but get %s, %v", res, err)
	}
}

func TestPostParser(t *testing.T) {
	input := `---
title: "Prove the Irrationality of Square Root of 2"
date: 2022-02-19T00:21:44+08:00
draft: false
tags: ["math", "number-theory"]
categories: ["Interest-Math"]
toc: true
summary: "Various methods of proving the irrationality of the square root of integer 2. Including the method by contradiction and the direct way."
---

# Introduction

The problem comes from the discussion about the proof of irrationality of $\sqrt{2}$ in the text books are usually by contradiction. In this essay, we discuss various methods for proving this problem. These proofs are mainly collected from the Internet. The origins would be specified as detailed as possible.
	`
	fmt.Printf("%+v\n", PostParser(input))
	t.Error("!")

}

func Test_contain(t *testing.T) {
	list := []string{"a", "bb", "ccc"}
	if contains(list, "aa") {
		t.Error("Failed for not_contain instance")
	}
	if !contains(list, "ccc") {
		t.Error("Failed for contain instance")
	}
}

func TestParseList(t *testing.T) {
	input1 := `["quantum","dynamics"]`
	input2 := `["Quantum-Mechanics"]`
	out1 := parseList(input1)
	for _, item := range out1 {
		if !contains([]string{"quantum", "dynamics"}, item) {
			t.Errorf("Failed for test instance 1 with output %s", out1)
		}
	}
	out2 := parseList(input2)
	for _, item := range out2 {
		if item != "Quantum-Mechanics" {
			t.Errorf("Failed for test instance 2 with output %s", out2)
		}
	}

}

func Test_extractWordsFromContent(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"01", args{"telskajet<div>abd</div>!!!test over"}, []string{"telskajet", "test"}},
		{"02", args{"test formula is $E=mc$, then"}, []string{"test", "formula"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractWordsFromContent(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractWordsFromContent() = [%v], want %v", strings.Join(got, ", "), tt.want)
			}
		})
	}
}
