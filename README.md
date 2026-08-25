# gomad
Embeddable Nomad Lisp interpreter written in Go.

<svg width="1143" height="1077" viewBox="0 0 1143 1077" fill="none" xmlns="http://www.w3.org/2000/svg">
<path fill-rule="evenodd" clip-rule="evenodd" d="M839.039 157.031C808.392 102.405 894.337 79.868 906.265 136.416L839.039 157.031Z" fill="#6AD7E5" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M1076.46 298.473C1133.91 308.517 1110.93 392.495 1058.25 366.729L1076.46 298.473Z" fill="#6AD7E5" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M765.977 623.477C769.245 637.353 768.876 669.06 748.24 663.977C724.571 662.934 741.386 632.288 733.829 618.155C744.842 616.483 756.545 616.619 765.977 623.477Z" fill="#F6D2A2" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path d="M748.24 663.976C750.125 658.279 753.909 653.215 753.782 646.977" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M594.706 488.455C583.087 479.95 568.347 483.364 555.349 479.586C542.732 476.395 545.931 467.15 546.482 465.175C545.808 463.371 544.821 464.566 544.834 460.728C549.071 441.837 572.485 448.346 586.021 450.393C599.066 459.05 596.458 475.258 594.706 488.455Z" fill="#F6D2A2" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path d="M546.483 465.175C551.274 461.151 557.72 461.694 563.111 459.632" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M855.799 143.829C853.149 135.229 853.064 127.061 863.011 123.996C872.219 121.158 877.483 127.73 880.133 136.331L855.799 143.829Z" fill="black"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M1065.25 341.963C1073.85 344.612 1082.01 344.698 1085.08 334.751C1087.92 325.542 1081.35 320.279 1072.74 317.628L1065.25 341.963Z" fill="black"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M1014.57 191.528C1057.26 234.212 1091.17 280.331 1071.02 342.199C1043.77 406.888 987.549 449.797 941.825 501.231C902.563 545.397 859.809 599.026 801.802 618.33C740.776 638.641 679.676 598.002 639.265 554.593C607.523 520.496 578.496 469.989 589.993 421.333C603.481 364.251 666.266 322.105 707.111 284.427C755.493 239.795 787.669 178.674 848.456 148.444C911.515 117.085 964.93 149.162 1014.57 191.528Z" fill="#6AD7E5" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M964.809 256.365C927.922 308.578 1012.52 369.62 1046.99 314.229C1077.91 264.559 1003.56 212.519 964.809 256.365Z" fill="white" stroke="black" stroke-width="2.9081" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M870.582 176.549C841.73 225.267 913.756 277.433 950.321 236.187C994.111 186.788 908.74 119.399 870.582 176.549Z" fill="white" stroke="black" stroke-width="2.8214" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M920.136 318.609C914.376 324.434 909.157 332.22 901.9 337.286C897.627 338.602 894.431 335.997 890.801 334.626C887.886 330.626 887.249 325.448 889.097 320.827C895.242 312.906 903.27 307.315 910.53 300.504L920.136 318.609Z" fill="white" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path d="M882.678 198.708C888.469 204.499 898.248 204.109 904.521 197.836C910.793 191.563 911.184 181.784 905.393 175.993C899.602 170.202 889.822 170.592 883.549 176.865C877.277 183.137 876.886 192.917 882.678 198.708Z" fill="black"/>
<path d="M893.355 197.569C894.72 198.935 897.119 198.749 898.713 197.155C900.307 195.561 900.493 193.162 899.128 191.797C897.762 190.431 895.363 190.617 893.769 192.211C892.175 193.805 891.99 196.204 893.355 197.569Z" fill="white"/>
<path d="M973.362 283.735C979.056 289.43 988.758 288.961 995.03 282.689C1001.3 276.416 1001.77 266.715 996.077 261.02C990.382 255.326 980.681 255.794 974.408 262.067C968.136 268.339 967.667 278.041 973.362 283.735Z" fill="black"/>
<path d="M983.96 282.518C985.303 283.86 987.683 283.656 989.277 282.062C990.872 280.468 991.075 278.088 989.733 276.745C988.39 275.402 986.009 275.606 984.415 277.2C982.821 278.794 982.617 281.175 983.96 282.518Z" fill="white"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M890.801 288.013C875.139 294.526 860.094 323.823 888.86 319.883C895.006 311.963 903.034 306.372 910.294 299.56L890.801 288.013Z" fill="white" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M913.72 268.226C904.466 260.426 887.357 263.595 883.77 276.049C879.018 292.544 903.686 294.085 910.765 303.331C919.966 312.917 917.923 330.354 933.639 329.727C951.12 329.029 947.968 308.173 942.544 294.833L913.72 268.226Z" fill="#F6D2A2" stroke="#231F20" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M913.09 266.709C927.231 251.305 957.243 277.623 948.204 293.315C939.185 308.968 903.529 279.043 913.09 266.709Z" fill="black"/>
<path d="M26.4293 538.13L545.939 26.0289L1065.45 538.13L545.939 1050.16L26.4293 538.13Z" fill="#6AD7E5"/>
<path d="M305.205 258.02L835.994 781.2L771.055 845.204L240.266 322.024L305.205 258.02Z" fill="black"/>
<path d="M516.306 581.044L581.245 645.049L358.331 864.759L293.392 800.755L516.306 581.044Z" fill="black"/>
<path d="M13.5485 524.715L532.324 13.3482L546.006 26.8298L27.2303 538.197L13.5485 524.715Z" fill="black"/>
<path d="M1065.05 538.464L532.39 13.4816L546.072 1.52588e-05L1078.73 524.982L1065.05 538.464Z" fill="black"/>
<path d="M545.939 1076.33L0 538.264L13.6818 524.782L559.621 1062.85L545.939 1076.33Z" fill="black"/>
<path d="M545.939 1049.7L1078.4 524.849L1092.08 538.33L559.621 1063.18L545.939 1049.7Z" fill="black"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M1012.81 475.799C1001.13 486.396 995.158 470.116 993.601 461.595C992.123 453.506 987.976 454.855 994.08 447.911C998.093 443.348 1001.89 438.566 1006.43 434.504C1014.32 441.082 1020.47 450.673 1022.09 460.912C1022.93 466.185 1021.57 478.998 1012.81 475.799Z" fill="#F6D2A2" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M1012.81 475.799C1011.79 473.577 1010.5 471.433 1010.87 468.87Z" fill="#C6B198"/>
<path d="M1012.81 475.799C1011.79 473.577 1010.5 471.433 1010.87 468.87" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M737.948 199.559C728.525 212.203 745.304 216.581 753.936 217.31C762.13 218.002 761.186 222.259 767.51 215.515C771.665 211.082 776.06 206.845 779.665 201.927C772.357 194.708 762.219 189.515 751.872 188.886C746.541 188.563 733.92 191.148 737.948 199.559Z" fill="#F6D2A2" stroke="black" stroke-width="3" stroke-linecap="round"/>
<path fill-rule="evenodd" clip-rule="evenodd" d="M737.948 199.559C740.257 200.357 742.516 201.433 745.032 200.822Z" fill="#C6B198"/>
<path d="M737.948 199.559C740.257 200.357 742.516 201.433 745.032 200.822" stroke="black" stroke-width="3" stroke-linecap="round"/>
</svg>

## What is this project about?
At its core, this is simply an implementation of [nomad](https://github.com/Moritisimor/nomad-lisp).

However, this implementation is more minimal and specifically made for embedding into go applications.

The difference to nomad is that gomad's API is much more suited for embedding and most of nomad's original I/O capabilities have been stripped from the prelude.

Also, be sure to check out [RobertFlexx's](https://github.com/robertflexx) nomad implementation written in rust, [romad](https://github.com/robertflexx/romad).

And if you want a BEAM port of nomad, be sure to check out **RobertFlexx's** [bomad](https://github.com/RobertFlexx/bomad)

## Getting started
Import the repository: 
```go
import "github.com/Moritisimor/gomad"
```

Then update your dependencies:
```bash
go mod tidy
```

And now you can use the `gomad` package within your project!

## Examples
### Instantiating the interpreter, registering natives and calling them through gomad
```go
package main

import (
	"fmt"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
)

func main() {
	interp := interpreter.New() // Create new interpreter
	interp.RegisterNative("hello", func(
        e []expr.Expression, 
        env *value.Env,
    ) (value.Value, error) {
		fmt.Println("Hello from Go!")
		return value.NewUnit(), nil // Unit is similar to null 
		// nil because no error occurred
	})

	evaluated, err := interp.DoString("(hello)") // S-Expressions
	if err != nil {
		fmt.Printf("Error while evaluating: %s\n", err.Error())
	} else {
		fmt.Printf("Evaluates to: %s\n", evaluated.String())
	}
}
```

### A simple logging function
```go
package main

import (
	"fmt"
	"log"

	"github.com/Moritisimor/gomad/expr"
	"github.com/Moritisimor/gomad/interpreter"
	"github.com/Moritisimor/gomad/value"
	"github.com/Moritisimor/gomad/eval"
)

func main() {
	interp := interpreter.New() // Create new interpreter
	interp.RegisterNative("log", func(
        e []expr.Expression, 
        env *value.Env,
    ) (value.Value, error) {
		if len(e) != 1 { // Check argument length
			return value.NewUnit(), fmt.Errorf("Error in call to log: Expected one argument, got %d", len(e))
		}

		// We expect a string as the argument, anything else is an error.
		logString, err := eval.GetString(e[0], env)
		if err != nil {
			return value.NewUnit(), fmt.Errorf("Error in call to log: %s", err.Error())
		}

		log.Println(logString)
		return value.NewUnit(), nil
	})

	interp.DoString("(log \"Gomad interpreter running!\")")
}
```
