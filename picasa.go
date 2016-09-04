package hello

type PicasaResponse struct {
  Feed PicasaFeed
}

type PicasaFeed struct {
  Entry []PicasaFeedEntry
}

type PicasaFeedEntry struct {
  Title PicasaValue
  Content PicasaContent
  OriginalVideo PicasaVideoInfo `json:"gphoto$originalvideo"`
  MediaGroup PicasaMediaGroup `json:"media$group"`
}

type PicasaValue struct {
  T string `json:"$t"`
}

type PicasaContent struct {
  Type string
  Src string
}

type PicasaVideoInfo struct {
  Type string
}

type PicasaMediaGroup struct {
  MediaContent []PicasaMediaContent `json:"media$content"`
}

type PicasaMediaContent struct {
  URL string
  Width int
  Height int
  Type string
  Medium string
}
