'use strict';

import {convert} from "./convert.js";

var slugify_cases = [
	{In: "", Exp: "-"},
	{In: "Hello  World 90", Exp: "hello-world-90"},
	{In: "Hello, 世界", Exp: "hello-世界"},
	{In: "90Things", Exp: "90things"},
	{In: "90 Things", Exp: "90-things"},
	{In: "KÜSIMUSED", Exp: "küsimused"},
	{In: "Küsimused Öösel", Exp: "küsimused-öösel"},
	{In: "nested / _paths", Exp: "nested/paths"},
	{In: "nested-/-paths", Exp: "nested/paths"},
	{In: "example_test.go", Exp: "example-test-go"},
	{In: "alpha + beta", Exp: "alpha-plus-beta"},
	{In: "alpha & beta", Exp: "alpha-amp-beta"},
	{In: "alpha # beta", Exp: "alpha-num-beta"},
	{In: "hello +/& world", Exp: "hello-plus/amp-world"},
	{In: "hello+/&world", Exp: "hello-plus/amp-world"},
	{In: "&Hello_世界/+!", Exp: "amp-hello-世界/plus-excl"}
];

for(var i = 0; i < slugify_cases.length; i += 1){
	var test = slugify_cases[i];
	var got = convert.Slugify(test.In);
	if(test.Exp != got) {
		console.log(i + " slugify '" + test.In + "': got '" + got + "' expected '" + test.Exp + "'");
	}
}
