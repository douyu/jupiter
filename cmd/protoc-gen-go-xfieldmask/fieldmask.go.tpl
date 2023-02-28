{{ range $idx, $m := .Message }}
    {{ if or $m.In $m.Out}}
        func (x *{{ $m.RequestName }}) fieldMaskWithMode(mode xfieldmask.MaskMode) *{{ $m.RequestName }}_FieldMask {
            fm := &{{ $m.RequestName }}_FieldMask{
                maskMode: mode,
                maskIn: xfieldmask.New(_fm_{{ $m.RequestName }}MaskIn),
                maskOut: xfieldmask.New(_fm_{{ $m.RequestName }}MaskOut),
            }

            return fm
        }

        // FieldMask_Filter generates *{{ $m.RequestName }} FieldMask with filter mode, so that
        // only the fields in the {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }} will be
        // appended into the {{ $m.RequestName }}.
        func (x *{{ $m.RequestName }}) FieldMaskFilter() *{{ $m.RequestName }}_FieldMask {
        	return x.fieldMaskWithMode(xfieldmask.MaskMode_Filter)
        }

        // FieldMask_Prune generates *{{ $m.RequestName }} FieldMask with prune mode, so that
        // only the fields NOT in the {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }} will be
        // appended into the {{ $m.RequestName }}.
        func (x *{{ $m.RequestName }}) FieldMaskPrune() *{{ $m.RequestName }}_FieldMask {
        	return x.fieldMaskWithMode(xfieldmask.MaskMode_Prune)
        }

        // {{ $m.RequestName }}_FieldMask provide helper functions to deal with FieldMask.
        type {{ $m.RequestName }}_FieldMask struct {
            maskMode xfieldmask.MaskMode
            maskIn xfieldmask.NestedFieldMask
            maskOut xfieldmask.NestedFieldMask
        }
    {{ end }}
    {{ if $m.In }}
        // _fm_{{ $m.RequestName }} is built in variable for {{ $m.RequestName }} to call FieldMask.Append
        var _fm_{{ $m.RequestName }} = new({{ $m.RequestName }})
        var _fm_{{ $m.RequestName }}MaskIn = make([]string, 0)

        {{ range $idx, $mf := $m.UpdateInFields }}
        {{ $maskFuncName := printf "MaskIn%s" $mf.UnderLineName }}
        {{ $maskedFuncName := printf "MaskedIn%s" $mf.UnderLineName }}
    
        // {{ $maskFuncName }} append {{ $mf.UnderLineName }} into {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }}.
        func (x *{{ $m.RequestName }}) {{ $maskFuncName }}() *{{ $m.RequestName }} {
            if x.{{ $m.IdentifyFieldGoName }} == nil {
                x.{{ $m.IdentifyFieldGoName }} = new(fieldmaskpb.FieldMask)
            }
            err := x.{{ $m.IdentifyFieldGoName }}.Append(_fm_{{ $m.RequestName }}, "{{ $mf.DotName }}")
            if err == nil{
                _fm_{{ $m.RequestName }}MaskIn = append(_fm_{{ $m.RequestName }}MaskIn,"{{ $mf.DotName }}")
            }
    
            return x
        }
    
        // {{ $maskedFuncName }} indicates the field {{ $mf.UnderLineName }} exists in the {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }} or not.
        func (x *{{ $m.RequestName }}_FieldMask) {{ $maskedFuncName }}() bool {
            if x.maskIn == nil {
                return true
            }
    
            return x.maskIn.Masked("{{ $mf.DotName }}")
        }
        {{ end }}
    {{ end }}
    {{ if $m.Out }}
        // _fm_{{ $m.ResponseName }} is built in variable for {{ $m.ResponseName }} to call FieldMask.Append
        var _fm_{{ $m.ResponseName }} = new({{ $m.ResponseName }})
        var _fm_{{ $m.RequestName }}MaskOut = make([]string, 0)

        {{ range $idx, $mf := $m.UpdateOutFields }}
        {{ $maskFuncName := printf "MaskOut%s" $mf.UnderLineName }}
        {{ $maskedFuncName := printf "MaskedOut%s" $mf.UnderLineName }}
    
        // {{ $maskFuncName }} append {{ $mf.UnderLineName }} into {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }}.
        func (x *{{ $m.RequestName }}) {{ $maskFuncName }}() *{{ $m.RequestName }} {
            if x.{{ $m.IdentifyFieldGoName }} == nil {
                x.{{ $m.IdentifyFieldGoName }} = new(fieldmaskpb.FieldMask)
            }
            err := x.{{ $m.IdentifyFieldGoName }}.Append(_fm_{{ $m.ResponseName }}, "{{ $mf.DotName }}")
            if err == nil{
                _fm_{{ $m.RequestName }}MaskOut = append(_fm_{{ $m.RequestName }}MaskOut, "{{ $mf.DotName }}")
            }
    
            return x
        }
    
        // {{ $maskedFuncName }} indicates the field {{ $mf.UnderLineName }} exists in the {{ $m.RequestName }}.{{ $m.IdentifyFieldGoName }} or not.
        func (x *{{ $m.RequestName }}_FieldMask) {{ $maskedFuncName }}() bool {
            if x.maskOut == nil {
                return true
            }
    
            return x.maskOut.Masked("{{ $mf.DotName }}")
        }
        {{ end }}
        // Mask only affects the fields in the {{ $m.RequestName }}.
        func (x *{{ $m.RequestName }}_FieldMask) Mask(m *{{$m.ResponseName}}) *{{$m.ResponseName}} {
           switch x.maskMode {
           case xfieldmask.MaskMode_Filter:
                x.maskOut.Filter(m)
           case xfieldmask.MaskMode_Prune:
                x.maskOut.Prune(m)
           }

           return m
        }
    {{ end }}
{{ end }}