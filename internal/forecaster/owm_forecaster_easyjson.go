// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package forecaster

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson1a67aff0DecodeWeatherReporterInternalForecaster(in *jlexer.Lexer, out *timedWeather) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "dt":
			out.Timestamp = int64(in.Int64())
		case "main":
			(out.Main).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1a67aff0EncodeWeatherReporterInternalForecaster(out *jwriter.Writer, in timedWeather) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"dt\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.Timestamp))
	}
	{
		const prefix string = ",\"main\":"
		out.RawString(prefix)
		(in.Main).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v timedWeather) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v timedWeather) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *timedWeather) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *timedWeather) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster(l, v)
}
func easyjson1a67aff0DecodeWeatherReporterInternalForecaster1(in *jlexer.Lexer, out *owmWeatherResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "main":
			(out.Main).UnmarshalEasyJSON(in)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1a67aff0EncodeWeatherReporterInternalForecaster1(out *jwriter.Writer, in owmWeatherResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"main\":"
		out.RawString(prefix[1:])
		(in.Main).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v owmWeatherResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v owmWeatherResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *owmWeatherResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *owmWeatherResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster1(l, v)
}
func easyjson1a67aff0DecodeWeatherReporterInternalForecaster2(in *jlexer.Lexer, out *owmForecastResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "list":
			if in.IsNull() {
				in.Skip()
				out.List = nil
			} else {
				in.Delim('[')
				if out.List == nil {
					if !in.IsDelim(']') {
						out.List = make([]timedWeather, 0, 4)
					} else {
						out.List = []timedWeather{}
					}
				} else {
					out.List = (out.List)[:0]
				}
				for !in.IsDelim(']') {
					var v1 timedWeather
					(v1).UnmarshalEasyJSON(in)
					out.List = append(out.List, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1a67aff0EncodeWeatherReporterInternalForecaster2(out *jwriter.Writer, in owmForecastResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"list\":"
		out.RawString(prefix[1:])
		if in.List == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.List {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v owmForecastResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v owmForecastResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *owmForecastResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *owmForecastResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster2(l, v)
}
func easyjson1a67aff0DecodeWeatherReporterInternalForecaster3(in *jlexer.Lexer, out *main) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "temp":
			out.Temperature = float64(in.Float64())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson1a67aff0EncodeWeatherReporterInternalForecaster3(out *jwriter.Writer, in main) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"temp\":"
		out.RawString(prefix[1:])
		out.Float64(float64(in.Temperature))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v main) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v main) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson1a67aff0EncodeWeatherReporterInternalForecaster3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *main) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *main) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson1a67aff0DecodeWeatherReporterInternalForecaster3(l, v)
}
