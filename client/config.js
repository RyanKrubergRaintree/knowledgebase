System.config({
	"baseURL": "/client/",
	"paths": {
		"*.js": "*.js",
		"*": "*.js",
		"vendor/*": "vendor/*"
	}
});

System.config({
	"map": {
		"jsx": "github:floatdrop/plugin-jsx@0.1.1",
		"react": "vendor/react@0.13.1/react.js",
		"ObjectId": "vendor/ObjectId.js",
	}
});
