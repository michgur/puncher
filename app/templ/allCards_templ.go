// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package templ

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/michgur/puncher/app/model"
import "github.com/michgur/puncher/app/templ/components"
import "fmt"

func openDrawer(path string) templ.ComponentScript {
	return templ.ComponentScript{
		Name: `__templ_openDrawer_5e51`,
		Function: `function __templ_openDrawer_5e51(path){document.querySelector('bottom-sheet').open();
	const otpInput = document.querySelector('otp-input');
	otpInput.setAttribute('path', path);
	otpInput.focus()
}`,
		Call:       templ.SafeScript(`__templ_openDrawer_5e51`, path),
		CallInline: templ.SafeScriptInline(`__templ_openDrawer_5e51`, path),
	}
}

func AllCards(cards []model.CardDetails) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = components.OTPInput().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = components.BottomSheet().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body class=\"bg-gray-100 flex flex-col h-dvh w-dvw text-gray-950 overflow-x-hidden\"><header class=\"flex p-4 gap-10\"><h1 class=\"basis-auto grow-0 shrink-0 text-3xl text-black font-semibold\">All Cards</h1><div class=\"flex-1 min-w-0 px-2 py-1 rounded-md bg-gray-300 flex gap-1\"><svg class=\"w-4 shrink-0 text-gray-600\" xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 512 512\"><path d=\"M416 208c0 45.9-14.9 88.3-40 122.7L502.6 457.4c12.5 12.5 12.5 32.8 0 45.3s-32.8 12.5-45.3 0L330.7 376c-34.4 25.2-76.8 40-122.7 40C93.1 416 0 322.9 0 208S93.1 0 208 0S416 93.1 416 208zM208 352a144 144 0 1 0 0-288 144 144 0 1 0 0 288z\" fill=\"CurrentColor\"></path></svg> <input class=\"flex-1 bg-transparent outline-none placeholder:text-gray-800\" type=\"text\" placeholder=\"Search Cards\"></div></header><div class=\"overflow-auto flex-1 flex flex-col gap-4 p-4 pt-0.5\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, card := range cards {
				templ_7745c5c3_Err = templ.RenderScriptItems(ctx, templ_7745c5c3_Buffer, openDrawer(fmt.Sprintf("/validate/%s", card.ID)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div onclick=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 templ.ComponentScript = openDrawer(fmt.Sprintf("/validate/%s", card.ID))
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ_7745c5c3_Var3.Call)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = Card(card, false).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div id=\"otp-sheet\"><h1>Punch a Slot</h1><p>Make a purchase at the business and ask the cashier for a puncher code</p><otp-input digits=\"4\"></otp-input></div><bottom-sheet content=\"#otp-sheet\"></bottom-sheet></body>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Index("All Cards").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
