package main

//
//func showCliOutput(){
//
//
//	for pkg, fInfo := range pkgInfo {
//
//		fmt.Sprintf("\n\n\n\n\n\n\n\n---------- Package : %v -------------------------------------------------------------------\n", pkg)
//
//		if showImports {
//			fmt.Printf("\n-- File Details -- \n ")
//
//			for _, v := range FileList {
//
//				fmt.Sprintf(`
//	%v --- %v
//	  Lines: %v, Size: %v,
//	  Funcs (Tot, Pub, Priv): (%v, %v, %v)
//	  Structs (Tot, Pub, Pri): (%v, %v, %v),
//	  Methods (Tot, Pub, Pri): (%v, %v, %v)
//	`, v.Filename, v.Name, v.NumberLines, v.Size, v.FunctionsCount, v.PublicFuncCount, v.PrivateFuncCount,
//					v.StructsCount, v.PublicStructCount, v.PrivateStructCount, v.MethodsCount, v.PublicMethodsCount, v.PrivateMethodsCount )
//			}
//		}
//
//		if showImports {
//			fmt.Printf("\n-- List of Imports -- \n%v \n ", strings.Join(fInfo.Imports, "\n"))
//		}
//		ii := 0
//		for _, v := range fInfo.Structs {
//
//			showStruct(v, fInfo, true, ii)
//			ii++
//		}
//
//		ii = 0
//		for _, v := range fInfo.Structs {
//
//			showStruct(v, fInfo, false, ii)
//			ii++
//		}
//
//
//		if showPublicFuncs || showFuncs {
//			fmt.Println("\n---   Public Functions ---------------------------------------------------")
//
//			for _, v := range fInfo.PublicFunctions {
//
//				fmt.Printf("  %v (%v)\t\t%v", v.Name, v.NumberLines, v.Content)
//			}
//		}
//
//		if showPrivateFuncs || showFuncs {
//			fmt.Println("\n---   Private Functions ---------------------------------------------------")
//
//			for _, v := range fInfo.PrivateFunctions {
//				fmt.Printf("  %v (%v)\t\t%v\n", v.Name, v.NumberLines, v.Content)
//			}
//		}
//	}
//
//
//}
//
//func showStruct(v FuncInfo, fInfo Funcs, publicLoop bool, ii int) {
//
//	if ((publicLoop && v.IsPublic) || (publicLoop == false && v.IsPublic == false)) && (showStructs || shortStructs) {
//
//		if ii == 0 && v.IsPublic {
//
//			fmt.Println("\n---   Public Structs ---------------------------------------------------")
//		}
//
//		if ii == 0 && v.IsPublic == false {
//
//			fmt.Println("\n---   Private Structs ---------------------------------------------------")
//		}
//
//		strCon := v.Content
//
//		if shortStructs {
//			strCon = strings.Split(v.Content, "\n")[0]
//		}
//
//		fmt.Printf("\n---- %v (%v) ---- \n %v \n", v.Name, v.NumberLines, strCon)
//
//		for i, vv := range fInfo.Methods[v.Name] {
//
//			if i == 0 && fInfo.LenPublicMethods(v.Name) > 0 {
//				fmt.Printf("\n\tPublic Methods: \n")
//			}
//
//			if vv.IsPublic {
//
//				fmt.Printf("\t  %v (%v)\t\t%v \n", vv.Name, vv.NumberLines, vv.Content)
//			}
//		}
//
//		for i, vv := range fInfo.Methods[v.Name] {
//
//			if i == 0 && fInfo.LenPrivateMethods(v.Name) > 0 {
//				fmt.Printf("\n\tPrivate Methods: \n")
//			}
//
//			if vv.IsPublic == false {
//
//				fmt.Printf("\t  %v (%v)\t\t%v \n", vv.Name, vv.NumberLines, vv.Content)
//			}
//		}
//
//		fmt.Printf("\n")
//	}
//}
