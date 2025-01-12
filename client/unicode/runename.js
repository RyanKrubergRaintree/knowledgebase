package("kb.unicode", function (exports) {
	"use strict";

	exports.RuneName = {
		"\u0021": "excl",
		"\u0022": "quot",
		"\u0023": "num",
		"\u0024": "dollar",
		"\u0025": "percnt",
		"\u0026": "amp",
		"\u0027": "apos",
		"\u0028": "lpar",
		"\u0029": "rpar",
		"\u002A": "ast",
		"\u002B": "plus",
		"\u002C": "comma",
		"\u002E": "period",
		"\u002F": "sol",
		"\u003A": "colon",
		"\u003B": "semi",
		"\u003C": "lt",
		"\u003D": "equals",
		"\u003E": "gt",
		"\u003F": "quest",
		"\u0040": "commat",
		"\u005B": "lsqb",
		"\u005C": "bsol",
		"\u005D": "rsqb",
		"\u005E": "hat",
		"\u005F": "lowbar",
		"\u0060": "grave",
		"\u007B": "lcub",
		"\u007C": "vert",
		"\u007D": "rcub",
		"\u00A1": "iexcl",
		"\u00A2": "cent",
		"\u00A3": "pound",
		"\u00A4": "curren",
		"\u00A5": "yen",
		"\u00A6": "brvbar",
		"\u00A7": "sect",
		"\u00A8": "uml",
		"\u00A9": "copy",
		"\u00AB": "laquo",
		"\u00AC": "not",
		"\u00AE": "reg",
		"\u00AF": "macr",
		"\u00B0": "deg",
		"\u00B1": "pm",
		"\u00B4": "acute",
		"\u00B6": "para",
		"\u00B7": "middot",
		"\u00B8": "cedil",
		"\u00BB": "raquo",
		"\u00BF": "iquest",
		"\u00D7": "times",
		"\u00F7": "div",
		"\u02D8": "breve",
		"\u02D9": "dot",
		"\u02DA": "ring",
		"\u02DB": "ogon",
		"\u02DC": "tilde",
		"\u02DD": "dblac",
		"\u03F6": "bepsi",
		"\u2010": "dash",
		"\u2013": "ndash",
		"\u2014": "mdash",
		"\u2015": "horbar",
		"\u2016": "vert",
		"\u2018": "lsquo",
		"\u2019": "rsquo",
		"\u201A": "sbquo",
		"\u201C": "ldquo",
		"\u201D": "rdquo",
		"\u201E": "bdquo",
		"\u2020": "dagger",
		"\u2021": "dagger",
		"\u2022": "bull",
		"\u2025": "nldr",
		"\u2026": "mldr",
		"\u2030": "permil",
		"\u2031": "pertenk",
		"\u2032": "prime",
		"\u2033": "prime",
		"\u2034": "tprime",
		"\u2035": "bprime",
		"\u2039": "lsaquo",
		"\u203A": "rsaquo",
		"\u203E": "oline",
		"\u2041": "caret",
		"\u2043": "hybull",
		"\u2044": "frasl",
		"\u204F": "bsemi",
		"\u2057": "qprime",
		"\u20AC": "euro",
		"\u2105": "incare",
		"\u2116": "numero",
		"\u2117": "copysr",
		"\u2118": "wp",
		"\u211E": "rx",
		"\u2122": "trade",
		"\u2127": "mho",
		"\u2129": "iiota",
		"\u2190": "larr",
		"\u2191": "uarr",
		"\u2192": "rarr",
		"\u2193": "darr",
		"\u2194": "harr",
		"\u2195": "varr",
		"\u2196": "nwarr",
		"\u2197": "nearr",
		"\u2198": "searr",
		"\u2199": "swarr",
		"\u219A": "nlarr",
		"\u219B": "nrarr",
		"\u219D": "rarrw",
		"\u219E": "larr",
		"\u219F": "uarr",
		"\u21A0": "rarr",
		"\u21A1": "darr",
		"\u21A2": "larrtl",
		"\u21A3": "rarrtl",
		"\u21A4": "mapstoleft",
		"\u21A5": "mapstoup",
		"\u21A6": "map",
		"\u21A7": "mapstodown",
		"\u21A9": "larrhk",
		"\u21AA": "rarrhk",
		"\u21AB": "larrlp",
		"\u21AC": "rarrlp",
		"\u21AD": "harrw",
		"\u21AE": "nharr",
		"\u21B0": "lsh",
		"\u21B1": "rsh",
		"\u21B2": "ldsh",
		"\u21B3": "rdsh",
		"\u21B5": "crarr",
		"\u21B6": "cularr",
		"\u21B7": "curarr",
		"\u21BA": "olarr",
		"\u21BB": "orarr",
		"\u21BC": "lharu",
		"\u21BD": "lhard",
		"\u21BE": "uharr",
		"\u21BF": "uharl",
		"\u21C0": "rharu",
		"\u21C1": "rhard",
		"\u21C2": "dharr",
		"\u21C3": "dharl",
		"\u21C4": "rlarr",
		"\u21C5": "udarr",
		"\u21C6": "lrarr",
		"\u21C7": "llarr",
		"\u21C8": "uuarr",
		"\u21C9": "rrarr",
		"\u21CA": "ddarr",
		"\u21CB": "lrhar",
		"\u21CC": "rlhar",
		"\u21CD": "nlarr",
		"\u21CE": "nharr",
		"\u21CF": "nrarr",
		"\u21D0": "larr",
		"\u21D1": "uarr",
		"\u21D2": "rarr",
		"\u21D3": "darr",
		"\u21D4": "iff",
		"\u21D5": "varr",
		"\u21D6": "nwarr",
		"\u21D7": "nearr",
		"\u21D8": "searr",
		"\u21D9": "swarr",
		"\u21DA": "laarr",
		"\u21DB": "raarr",
		"\u21DD": "zigrarr",
		"\u21E4": "larrb",
		"\u21E5": "rarrb",
		"\u21F5": "duarr",
		"\u21FD": "loarr",
		"\u21FE": "roarr",
		"\u21FF": "hoarr",
		"\u2200": "forall",
		"\u2201": "comp",
		"\u2202": "part",
		"\u2203": "exist",
		"\u2204": "nexist",
		"\u2205": "empty",
		"\u2207": "del",
		"\u2208": "in",
		"\u2209": "notin",
		"\u220B": "ni",
		"\u220C": "notni",
		"\u220F": "prod",
		"\u2210": "coprod",
		"\u2211": "sum",
		"\u2212": "minus",
		"\u2213": "mp",
		"\u2214": "plusdo",
		"\u2216": "setmn",
		"\u2217": "lowast",
		"\u2218": "compfn",
		"\u221A": "sqrt",
		"\u221D": "prop",
		"\u221E": "infin",
		"\u221F": "angrt",
		"\u2220": "ang",
		"\u2221": "angmsd",
		"\u2222": "angsph",
		"\u2223": "mid",
		"\u2224": "nmid",
		"\u2225": "par",
		"\u2226": "npar",
		"\u2227": "and",
		"\u2228": "or",
		"\u2229": "cap",
		"\u222A": "cup",
		"\u222B": "int",
		"\u222C": "int",
		"\u222D": "tint",
		"\u222E": "oint",
		"\u222F": "conint",
		"\u2230": "cconint",
		"\u2231": "cwint",
		"\u2232": "cwconint",
		"\u2233": "awconint",
		"\u2234": "there4",
		"\u2235": "becaus",
		"\u2236": "ratio",
		"\u2237": "colon",
		"\u2238": "minusd",
		"\u223A": "mddot",
		"\u223B": "homtht",
		"\u223C": "sim",
		"\u223D": "bsim",
		"\u223E": "ac",
		"\u223F": "acd",
		"\u2240": "wr",
		"\u2241": "nsim",
		"\u2242": "esim",
		"\u2243": "sime",
		"\u2244": "nsime",
		"\u2245": "cong",
		"\u2246": "simne",
		"\u2247": "ncong",
		"\u2248": "ap",
		"\u2249": "nap",
		"\u224A": "ape",
		"\u224B": "apid",
		"\u224C": "bcong",
		"\u224D": "cupcap",
		"\u224E": "bump",
		"\u224F": "bumpe",
		"\u2250": "doteq",
		"\u2251": "edot",
		"\u2252": "efdot",
		"\u2253": "erdot",
		"\u2254": "assign",
		"\u2255": "ecolon",
		"\u2256": "ecir",
		"\u2257": "cire",
		"\u2259": "wedgeq",
		"\u225A": "veeeq",
		"\u225C": "trie",
		"\u225F": "equest",
		"\u2260": "ne",
		"\u2261": "equiv",
		"\u2262": "nequiv",
		"\u2264": "le",
		"\u2265": "ge",
		"\u2266": "le",
		"\u2267": "ge",
		"\u2268": "lne",
		"\u2269": "gne",
		"\u226A": "lt",
		"\u226B": "gt",
		"\u226C": "twixt",
		"\u226D": "notcupcap",
		"\u226E": "nlt",
		"\u226F": "ngt",
		"\u2270": "nle",
		"\u2271": "nge",
		"\u2272": "lsim",
		"\u2273": "gsim",
		"\u2274": "nlsim",
		"\u2275": "ngsim",
		"\u2276": "lg",
		"\u2277": "gl",
		"\u2278": "ntlg",
		"\u2279": "ntgl",
		"\u227A": "pr",
		"\u227B": "sc",
		"\u227C": "prcue",
		"\u227D": "sccue",
		"\u227E": "prsim",
		"\u227F": "scsim",
		"\u2280": "npr",
		"\u2281": "nsc",
		"\u2282": "sub",
		"\u2283": "sup",
		"\u2284": "nsub",
		"\u2285": "nsup",
		"\u2286": "sube",
		"\u2287": "supe",
		"\u2288": "nsube",
		"\u2289": "nsupe",
		"\u228A": "subne",
		"\u228B": "supne",
		"\u228D": "cupdot",
		"\u228E": "uplus",
		"\u228F": "sqsub",
		"\u2290": "sqsup",
		"\u2291": "sqsube",
		"\u2292": "sqsupe",
		"\u2293": "sqcap",
		"\u2294": "sqcup",
		"\u2295": "oplus",
		"\u2296": "ominus",
		"\u2297": "otimes",
		"\u2298": "osol",
		"\u2299": "odot",
		"\u229A": "ocir",
		"\u229B": "oast",
		"\u229D": "odash",
		"\u229E": "plusb",
		"\u229F": "minusb",
		"\u22A0": "timesb",
		"\u22A1": "sdotb",
		"\u22A2": "vdash",
		"\u22A3": "dashv",
		"\u22A4": "top",
		"\u22A5": "bot",
		"\u22A7": "models",
		"\u22A8": "vdash",
		"\u22A9": "vdash",
		"\u22AA": "vvdash",
		"\u22AB": "vdash",
		"\u22AC": "nvdash",
		"\u22AD": "nvdash",
		"\u22AE": "nvdash",
		"\u22AF": "nvdash",
		"\u22B0": "prurel",
		"\u22B2": "vltri",
		"\u22B3": "vrtri",
		"\u22B4": "ltrie",
		"\u22B5": "rtrie",
		"\u22B6": "origof",
		"\u22B7": "imof",
		"\u22B8": "mumap",
		"\u22B9": "hercon",
		"\u22BA": "intcal",
		"\u22BB": "veebar",
		"\u22BD": "barvee",
		"\u22BE": "angrtvb",
		"\u22BF": "lrtri",
		"\u22C0": "wedge",
		"\u22C1": "vee",
		"\u22C2": "xcap",
		"\u22C3": "xcup",
		"\u22C4": "diam",
		"\u22C5": "sdot",
		"\u22C6": "star",
		"\u22C7": "divonx",
		"\u22C8": "bowtie",
		"\u22C9": "ltimes",
		"\u22CA": "rtimes",
		"\u22CB": "lthree",
		"\u22CC": "rthree",
		"\u22CD": "bsime",
		"\u22CE": "cuvee",
		"\u22CF": "cuwed",
		"\u22D0": "sub",
		"\u22D1": "sup",
		"\u22D2": "cap",
		"\u22D3": "cup",
		"\u22D4": "fork",
		"\u22D5": "epar",
		"\u22D6": "ltdot",
		"\u22D7": "gtdot",
		"\u22D8": "ll",
		"\u22D9": "gg",
		"\u22DA": "leg",
		"\u22DB": "gel",
		"\u22DE": "cuepr",
		"\u22DF": "cuesc",
		"\u22E0": "nprcue",
		"\u22E1": "nsccue",
		"\u22E2": "nsqsube",
		"\u22E3": "nsqsupe",
		"\u22E6": "lnsim",
		"\u22E7": "gnsim",
		"\u22E8": "prnsim",
		"\u22E9": "scnsim",
		"\u22EA": "nltri",
		"\u22EB": "nrtri",
		"\u22EC": "nltrie",
		"\u22ED": "nrtrie",
		"\u22EE": "vellip",
		"\u22EF": "ctdot",
		"\u22F0": "utdot",
		"\u22F1": "dtdot",
		"\u22F2": "disin",
		"\u22F3": "isinsv",
		"\u22F4": "isins",
		"\u22F5": "isindot",
		"\u22F6": "notinvc",
		"\u22F7": "notinvb",
		"\u22F9": "isine",
		"\u22FA": "nisd",
		"\u22FB": "xnis",
		"\u22FC": "nis",
		"\u22FD": "notnivc",
		"\u22FE": "notnivb",
		"\u2305": "barwed",
		"\u2306": "barwed",
		"\u2308": "lceil",
		"\u2309": "rceil",
		"\u230A": "lfloor",
		"\u230B": "rfloor",
		"\u230C": "drcrop",
		"\u230D": "dlcrop",
		"\u230E": "urcrop",
		"\u230F": "ulcrop",
		"\u2310": "bnot",
		"\u2312": "profline",
		"\u2313": "profsurf",
		"\u2315": "telrec",
		"\u2316": "target",
		"\u231C": "ulcorn",
		"\u231D": "urcorn",
		"\u231E": "dlcorn",
		"\u231F": "drcorn",
		"\u2322": "frown",
		"\u2323": "smile",
		"\u232D": "cylcty",
		"\u232E": "profalar",
		"\u2336": "topbot",
		"\u233D": "ovbar",
		"\u233F": "solbar",
		"\u237C": "angzarr",
		"\u23B0": "lmoust",
		"\u23B1": "rmoust",
		"\u23B4": "tbrk",
		"\u23B5": "bbrk",
		"\u23B6": "bbrktbrk",
		"\u23DC": "overparenthesis",
		"\u23DD": "underparenthesis",
		"\u23DE": "overbrace",
		"\u23DF": "underbrace",
		"\u23E2": "trpezium",
		"\u23E7": "elinters",
		"\u2423": "blank",
		"\u24C8": "os",
		"\u2500": "boxh",
		"\u2502": "boxv",
		"\u250C": "boxdr",
		"\u2510": "boxdl",
		"\u2514": "boxur",
		"\u2518": "boxul",
		"\u251C": "boxvr",
		"\u2524": "boxvl",
		"\u252C": "boxhd",
		"\u2534": "boxhu",
		"\u253C": "boxvh",
		"\u2550": "boxh",
		"\u2551": "boxv",
		"\u2552": "boxdr",
		"\u2553": "boxdr",
		"\u2554": "boxdr",
		"\u2555": "boxdl",
		"\u2556": "boxdl",
		"\u2557": "boxdl",
		"\u2558": "boxur",
		"\u2559": "boxur",
		"\u255A": "boxur",
		"\u255B": "boxul",
		"\u255C": "boxul",
		"\u255D": "boxul",
		"\u255E": "boxvr",
		"\u255F": "boxvr",
		"\u2560": "boxvr",
		"\u2561": "boxvl",
		"\u2562": "boxvl",
		"\u2563": "boxvl",
		"\u2564": "boxhd",
		"\u2565": "boxhd",
		"\u2566": "boxhd",
		"\u2567": "boxhu",
		"\u2568": "boxhu",
		"\u2569": "boxhu",
		"\u256A": "boxvh",
		"\u256B": "boxvh",
		"\u256C": "boxvh",
		"\u2580": "uhblk",
		"\u2584": "lhblk",
		"\u2588": "block",
		"\u2591": "blk14",
		"\u2592": "blk12",
		"\u2593": "blk34",
		"\u25A1": "squ",
		"\u25AA": "squf",
		"\u25AB": "emptyverysmallsquare",
		"\u25AD": "rect",
		"\u25AE": "marker",
		"\u25B1": "fltns",
		"\u25B3": "xutri",
		"\u25B4": "utrif",
		"\u25B5": "utri",
		"\u25B8": "rtrif",
		"\u25B9": "rtri",
		"\u25BD": "xdtri",
		"\u25BE": "dtrif",
		"\u25BF": "dtri",
		"\u25C2": "ltrif",
		"\u25C3": "ltri",
		"\u25CA": "loz",
		"\u25CB": "cir",
		"\u25EC": "tridot",
		"\u25EF": "xcirc",
		"\u25F8": "ultri",
		"\u25F9": "urtri",
		"\u25FA": "lltri",
		"\u25FB": "emptysmallsquare",
		"\u25FC": "filledsmallsquare",
		"\u2605": "starf",
		"\u2606": "star",
		"\u260E": "phone",
		"\u2640": "female",
		"\u2642": "male",
		"\u2660": "spades",
		"\u2663": "clubs",
		"\u2665": "hearts",
		"\u2666": "diams",
		"\u266A": "sung",
		"\u266D": "flat",
		"\u266E": "natur",
		"\u266F": "sharp",
		"\u2713": "check",
		"\u2717": "cross",
		"\u2720": "malt",
		"\u2736": "sext",
		"\u2758": "verticalseparator",
		"\u2772": "lbbrk",
		"\u2773": "rbbrk",
		"\u27C8": "bsolhsub",
		"\u27C9": "suphsol",
		"\u27E6": "lobrk",
		"\u27E7": "robrk",
		"\u27E8": "lang",
		"\u27E9": "rang",
		"\u27EA": "lang",
		"\u27EB": "rang",
		"\u27EC": "loang",
		"\u27ED": "roang",
		"\u27F5": "xlarr",
		"\u27F6": "xrarr",
		"\u27F7": "xharr",
		"\u27F8": "xlarr",
		"\u27F9": "xrarr",
		"\u27FA": "xharr",
		"\u27FC": "xmap",
		"\u27FF": "dzigrarr",
		"\u2902": "nvlarr",
		"\u2903": "nvrarr",
		"\u2904": "nvharr",
		"\u2905": "map",
		"\u290C": "lbarr",
		"\u290D": "rbarr",
		"\u290E": "lbarr",
		"\u290F": "rbarr",
		"\u2910": "rbarr",
		"\u2911": "ddotrahd",
		"\u2912": "uparrowbar",
		"\u2913": "downarrowbar",
		"\u2916": "rarrtl",
		"\u2919": "latail",
		"\u291A": "ratail",
		"\u291B": "latail",
		"\u291C": "ratail",
		"\u291D": "larrfs",
		"\u291E": "rarrfs",
		"\u291F": "larrbfs",
		"\u2920": "rarrbfs",
		"\u2923": "nwarhk",
		"\u2924": "nearhk",
		"\u2925": "searhk",
		"\u2926": "swarhk",
		"\u2927": "nwnear",
		"\u2928": "toea",
		"\u2929": "tosa",
		"\u292A": "swnwar",
		"\u2933": "rarrc",
		"\u2935": "cudarrr",
		"\u2936": "ldca",
		"\u2937": "rdca",
		"\u2938": "cudarrl",
		"\u2939": "larrpl",
		"\u293C": "curarrm",
		"\u293D": "cularrp",
		"\u2945": "rarrpl",
		"\u2948": "harrcir",
		"\u2949": "uarrocir",
		"\u294A": "lurdshar",
		"\u294B": "ldrushar",
		"\u294E": "leftrightvector",
		"\u294F": "rightupdownvector",
		"\u2950": "downleftrightvector",
		"\u2951": "leftupdownvector",
		"\u2952": "leftvectorbar",
		"\u2953": "rightvectorbar",
		"\u2954": "rightupvectorbar",
		"\u2955": "rightdownvectorbar",
		"\u2956": "downleftvectorbar",
		"\u2957": "downrightvectorbar",
		"\u2958": "leftupvectorbar",
		"\u2959": "leftdownvectorbar",
		"\u295A": "leftteevector",
		"\u295B": "rightteevector",
		"\u295C": "rightupteevector",
		"\u295D": "rightdownteevector",
		"\u295E": "downleftteevector",
		"\u295F": "downrightteevector",
		"\u2960": "leftupteevector",
		"\u2961": "leftdownteevector",
		"\u2962": "lhar",
		"\u2963": "uhar",
		"\u2964": "rhar",
		"\u2965": "dhar",
		"\u2966": "luruhar",
		"\u2967": "ldrdhar",
		"\u2968": "ruluhar",
		"\u2969": "rdldhar",
		"\u296A": "lharul",
		"\u296B": "llhard",
		"\u296C": "rharul",
		"\u296D": "lrhard",
		"\u296E": "udhar",
		"\u296F": "duhar",
		"\u2970": "roundimplies",
		"\u2971": "erarr",
		"\u2972": "simrarr",
		"\u2973": "larrsim",
		"\u2974": "rarrsim",
		"\u2975": "rarrap",
		"\u2976": "ltlarr",
		"\u2978": "gtrarr",
		"\u2979": "subrarr",
		"\u297B": "suplarr",
		"\u297C": "lfisht",
		"\u297D": "rfisht",
		"\u297E": "ufisht",
		"\u297F": "dfisht",
		"\u2985": "lopar",
		"\u2986": "ropar",
		"\u298B": "lbrke",
		"\u298C": "rbrke",
		"\u298D": "lbrkslu",
		"\u298E": "rbrksld",
		"\u298F": "lbrksld",
		"\u2990": "rbrkslu",
		"\u2991": "langd",
		"\u2992": "rangd",
		"\u2993": "lparlt",
		"\u2994": "rpargt",
		"\u2995": "gtlpar",
		"\u2996": "ltrpar",
		"\u299A": "vzigzag",
		"\u299C": "vangrt",
		"\u299D": "angrtvbd",
		"\u29A4": "ange",
		"\u29A5": "range",
		"\u29A6": "dwangle",
		"\u29A7": "uwangle",
		"\u29A8": "angmsdaa",
		"\u29A9": "angmsdab",
		"\u29AA": "angmsdac",
		"\u29AB": "angmsdad",
		"\u29AC": "angmsdae",
		"\u29AD": "angmsdaf",
		"\u29AE": "angmsdag",
		"\u29AF": "angmsdah",
		"\u29B0": "bemptyv",
		"\u29B1": "demptyv",
		"\u29B2": "cemptyv",
		"\u29B3": "raemptyv",
		"\u29B4": "laemptyv",
		"\u29B5": "ohbar",
		"\u29B6": "omid",
		"\u29B7": "opar",
		"\u29B9": "operp",
		"\u29BB": "olcross",
		"\u29BC": "odsold",
		"\u29BE": "olcir",
		"\u29BF": "ofcir",
		"\u29C0": "olt",
		"\u29C1": "ogt",
		"\u29C2": "cirscir",
		"\u29C3": "cire",
		"\u29C4": "solb",
		"\u29C5": "bsolb",
		"\u29C9": "boxbox",
		"\u29CD": "trisb",
		"\u29CE": "rtriltri",
		"\u29CF": "lefttrianglebar",
		"\u29D0": "righttrianglebar",
		"\u29DC": "iinfin",
		"\u29DD": "infintie",
		"\u29DE": "nvinfin",
		"\u29E3": "eparsl",
		"\u29E4": "smeparsl",
		"\u29E5": "eqvparsl",
		"\u29EB": "lozf",
		"\u29F4": "ruledelayed",
		"\u29F6": "dsol",
		"\u2A00": "xodot",
		"\u2A01": "xoplus",
		"\u2A02": "xotime",
		"\u2A04": "xuplus",
		"\u2A06": "xsqcup",
		"\u2A0C": "qint",
		"\u2A0D": "fpartint",
		"\u2A10": "cirfnint",
		"\u2A11": "awint",
		"\u2A12": "rppolint",
		"\u2A13": "scpolint",
		"\u2A14": "npolint",
		"\u2A15": "pointint",
		"\u2A16": "quatint",
		"\u2A17": "intlarhk",
		"\u2A22": "pluscir",
		"\u2A23": "plusacir",
		"\u2A24": "simplus",
		"\u2A25": "plusdu",
		"\u2A26": "plussim",
		"\u2A27": "plustwo",
		"\u2A29": "mcomma",
		"\u2A2A": "minusdu",
		"\u2A2D": "loplus",
		"\u2A2E": "roplus",
		"\u2A2F": "cross",
		"\u2A30": "timesd",
		"\u2A31": "timesbar",
		"\u2A33": "smashp",
		"\u2A34": "lotimes",
		"\u2A35": "rotimes",
		"\u2A36": "otimesas",
		"\u2A37": "otimes",
		"\u2A38": "odiv",
		"\u2A39": "triplus",
		"\u2A3A": "triminus",
		"\u2A3B": "tritime",
		"\u2A3C": "iprod",
		"\u2A3F": "amalg",
		"\u2A40": "capdot",
		"\u2A42": "ncup",
		"\u2A43": "ncap",
		"\u2A44": "capand",
		"\u2A45": "cupor",
		"\u2A46": "cupcap",
		"\u2A47": "capcup",
		"\u2A48": "cupbrcap",
		"\u2A49": "capbrcup",
		"\u2A4A": "cupcup",
		"\u2A4B": "capcap",
		"\u2A4C": "ccups",
		"\u2A4D": "ccaps",
		"\u2A50": "ccupssm",
		"\u2A53": "and",
		"\u2A54": "or",
		"\u2A55": "andand",
		"\u2A56": "oror",
		"\u2A57": "orslope",
		"\u2A58": "andslope",
		"\u2A5A": "andv",
		"\u2A5B": "orv",
		"\u2A5C": "andd",
		"\u2A5D": "ord",
		"\u2A5F": "wedbar",
		"\u2A66": "sdote",
		"\u2A6A": "simdot",
		"\u2A6D": "congdot",
		"\u2A6E": "easter",
		"\u2A6F": "apacir",
		"\u2A70": "ape",
		"\u2A71": "eplus",
		"\u2A72": "pluse",
		"\u2A73": "esim",
		"\u2A74": "colone",
		"\u2A75": "equal",
		"\u2A77": "eddot",
		"\u2A78": "equivdd",
		"\u2A79": "ltcir",
		"\u2A7A": "gtcir",
		"\u2A7B": "ltquest",
		"\u2A7C": "gtquest",
		"\u2A7D": "les",
		"\u2A7E": "ges",
		"\u2A7F": "lesdot",
		"\u2A80": "gesdot",
		"\u2A81": "lesdoto",
		"\u2A82": "gesdoto",
		"\u2A83": "lesdotor",
		"\u2A84": "gesdotol",
		"\u2A85": "lap",
		"\u2A86": "gap",
		"\u2A87": "lne",
		"\u2A88": "gne",
		"\u2A89": "lnap",
		"\u2A8A": "gnap",
		"\u2A8B": "leg",
		"\u2A8C": "gel",
		"\u2A8D": "lsime",
		"\u2A8E": "gsime",
		"\u2A8F": "lsimg",
		"\u2A90": "gsiml",
		"\u2A91": "lge",
		"\u2A92": "gle",
		"\u2A93": "lesges",
		"\u2A94": "gesles",
		"\u2A95": "els",
		"\u2A96": "egs",
		"\u2A97": "elsdot",
		"\u2A98": "egsdot",
		"\u2A99": "el",
		"\u2A9A": "eg",
		"\u2A9D": "siml",
		"\u2A9E": "simg",
		"\u2A9F": "simle",
		"\u2AA0": "simge",
		"\u2AA1": "lessless",
		"\u2AA2": "greatergreater",
		"\u2AA4": "glj",
		"\u2AA5": "gla",
		"\u2AA6": "ltcc",
		"\u2AA7": "gtcc",
		"\u2AA8": "lescc",
		"\u2AA9": "gescc",
		"\u2AAA": "smt",
		"\u2AAB": "lat",
		"\u2AAC": "smte",
		"\u2AAD": "late",
		"\u2AAE": "bumpe",
		"\u2AAF": "pre",
		"\u2AB0": "sce",
		"\u2AB3": "pre",
		"\u2AB4": "sce",
		"\u2AB5": "prne",
		"\u2AB6": "scne",
		"\u2AB7": "prap",
		"\u2AB8": "scap",
		"\u2AB9": "prnap",
		"\u2ABA": "scnap",
		"\u2ABB": "pr",
		"\u2ABC": "sc",
		"\u2ABD": "subdot",
		"\u2ABE": "supdot",
		"\u2ABF": "subplus",
		"\u2AC0": "supplus",
		"\u2AC1": "submult",
		"\u2AC2": "supmult",
		"\u2AC3": "subedot",
		"\u2AC4": "supedot",
		"\u2AC5": "sube",
		"\u2AC6": "supe",
		"\u2AC7": "subsim",
		"\u2AC8": "supsim",
		"\u2ACB": "subne",
		"\u2ACC": "supne",
		"\u2ACF": "csub",
		"\u2AD0": "csup",
		"\u2AD1": "csube",
		"\u2AD2": "csupe",
		"\u2AD3": "subsup",
		"\u2AD4": "supsub",
		"\u2AD5": "subsub",
		"\u2AD6": "supsup",
		"\u2AD7": "suphsub",
		"\u2AD8": "supdsub",
		"\u2AD9": "forkv",
		"\u2ADA": "topfork",
		"\u2ADB": "mlcp",
		"\u2AE4": "dashv",
		"\u2AE6": "vdashl",
		"\u2AE7": "barv",
		"\u2AE8": "vbar",
		"\u2AE9": "vbarv",
		"\u2AEB": "vbar",
		"\u2AEC": "not",
		"\u2AED": "bnot",
		"\u2AEE": "rnmid",
		"\u2AEF": "cirmid",
		"\u2AF0": "midcir",
		"\u2AF1": "topcir",
		"\u2AF2": "nhpar",
		"\u2AF3": "parsim",
		"\u2AFD": "parsl"
	};
});
