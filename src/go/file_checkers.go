package main

import (
	"bytes"
	"fmt"
	"github.com/alfg/mp4"
)

func check_ogg(data []byte) string {
	reader, header, err := NewWith(bytes.NewReader(data)) // TODO: have oggreader in its package
	if reader != nil && err == nil && header != nil && header.SampleRate == 16000 && header.Channels == 1 {
		return ""
	} else {
		return fmt.Sprintf("An ogg audio file was provided, but it has not the correct format.\nNeed opus, channel: 1, rate: 16000, error nil.\nGot channels: %d , rate: %d, error %s.", header.Channels, header.SampleRate, err)
	}
}

func check_mp4(data []byte) string {
	mp4, err := mp4.OpenFromBytes(data)
	if err != nil {
		return fmt.Sprintf("An mp4 video file was provided, but it could not be analyzed because of %v.", err)
	}
	// NOTE: I assume brand must be mp42
	// TODO: also accept compatible brands with slices.Contains(mp4.Ftyp.CompatibleBrands, "mp42")
	if mp4.Ftyp.MajorBrand == "mp42" && mp4.Ftyp.MinorVersion == 0 && mp4.Moov.IsFragmented == false {
		// NOTE: I assume video track must be at 0
		// presence of Avc1 box (Box.Name "avc1") indicates h264
		// see https://developer.apple.com/library/archive/documentation/QuickTime/QTFF/QTFFChap3/qtff3.html
		video_ok := false
		audio_ok := true // TODO: actually check audio
		if len(mp4.Moov.Traks) == 1 || len(mp4.Moov.Traks) == 2 {
			trak := mp4.Moov.Traks[0]
			video_ok = trak.Mdia.Minf.Stbl != nil && trak.Mdia.Minf.Stbl.Stsd.Avc1 != nil
			// TODO: access h264 stream and check for yuv420
		}
		// TODO: assert mp4.Edts == nil ?
		if video_ok && audio_ok {
			return ""
		} else {
			return fmt.Sprintf("An mp42 video file was provided, but it has not the correct format.\nFirst track avc1 (aka. h264) video. Second track mp4a (aka. aac) audio (optional).")
		}
	} else {
		return fmt.Sprintf("An mp4 video file was provided, but it has not the correct container.\nNeed brand: mp42, version: 0, fragmented: false.")
	}
	/*
	   handler.log.Infof("mp4: %+v", mp4)
	   handler.log.Infof("mp4.Ftyp: %+v", mp4.Ftyp)
	   handler.log.Infof("mp4.Moov: %+v", mp4.Moov)
	   handler.log.Infof("mp4.Moov.Mvhd: %+v", mp4.Moov.Mvhd)
	   for i, trak := range mp4.Moov.Traks {
	       handler.log.Infof("trak[%d]: %+v", i, trak)
	       handler.log.Infof("trak[%d].Tkhd: %+v", i, trak.Tkhd)
	       handler.log.Infof("trak[%d].Mdia: %+v", i, trak.Mdia)
	       handler.log.Infof("trak[%d].Mdia.Hdlr: %+v", i, trak.Mdia.Hdlr) // Handler:"vide", Handler:"soun"
	       handler.log.Infof("trak[%d].Mdia.Mdhd: %+v", i, trak.Mdia.Mdhd)
	       handler.log.Infof("trak[%d].Mdia.Minf: %+v", i, trak.Mdia.Minf)
	       if trak.Mdia.Minf.Vmhd != nil {
	           handler.log.Infof("trak[%d].Mdia.Minf.Vmhd: %+v", i, trak.Mdia.Minf.Vmhd)
	       }
	       if trak.Mdia.Minf.Hmhd != nil {
	           handler.log.Infof("trak[%d].Mdia.Minf.Hmhd: %+v", i, trak.Mdia.Minf.Hmhd)
	       }
	       if trak.Mdia.Minf.Stbl != nil {
	           handler.log.Infof("trak[%d].Mdia.Minf.Stbl: %+v", i, trak.Mdia.Minf.Stbl)
	           handler.log.Infof("trak[%d].Mdia.Minf.Stbl.Stts: %+v", i, trak.Mdia.Minf.Stbl.Stts)
	           handler.log.Infof("trak[%d].Mdia.Minf.Stbl.Stsd: %+v", i, trak.Mdia.Minf.Stbl.Stsd)
	           if trak.Mdia.Minf.Stbl.Stsd.Avc1 != nil {
	               handler.log.Infof("trak[%d].Mdia.Minf.Stbl.Stsd.Avc1: %+v", i, trak.Mdia.Minf.Stbl.Stsd.Avc1)
	               handler.log.Infof("trak[%d].Mdia.Minf.Stbl.Stsd.Avc1.Box: %+v", i, trak.Mdia.Minf.Stbl.Stsd.Avc1.Box)
	           }
	       }
	   }
	*/
}
