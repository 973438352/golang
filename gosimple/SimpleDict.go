package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type CaiyunDictRequest struct {
	Trans_type string `json:"trans_type"`
	Source     string `json:"source"`
	UserId     string `json:"user_id"`
}

type CaiyunDictResponse struct {
	Rc   int `json:"rc"`
	Wiki struct {
		KnownInLaguages int `json:"known_in_laguages"`
		Description     struct {
			Source string      `json:"source"`
			Target interface{} `json:"target"`
		} `json:"description"`
		ID   string `json:"id"`
		Item struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"item"`
		ImageURL  string `json:"image_url"`
		IsSubject string `json:"is_subject"`
		Sitelink  string `json:"sitelink"`
	} `json:"wiki"`
	Dictionary struct {
		Prons struct {
			EnUs string `json:"en-us"`
			En   string `json:"en"`
		} `json:"prons"`
		Explanations []string      `json:"explanations"`
		Synonym      []string      `json:"synonym"`
		Antonym      []string      `json:"antonym"`
		WqxExample   [][]string    `json:"wqx_example"`
		Entry        string        `json:"entry"`
		Type         string        `json:"type"`
		Related      []interface{} `json:"related"`
		Source       string        `json:"source"`
	} `json:"dictionary"`
}

func queryCaiyun(word string, wg *sync.WaitGroup) {
	client := &http.Client{}
	//var data = strings.NewReader(`{"trans_type":"en2zh","source":"good"}`)
	request := CaiyunDictRequest{Trans_type: "en2zh", Source: word}
	buf, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
	}
	var data = bytes.NewReader(buf)
	req, err := http.NewRequest("POST", "https://api.interpreter.caiyunai.com/v1/dict", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "api.interpreter.caiyunai.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("app-name", "xy")
	req.Header.Set("content-type", "application/json;charset=UTF-8")
	req.Header.Set("device-id", "")
	req.Header.Set("origin", "https://fanyi.caiyunapp.com")
	req.Header.Set("os-type", "web")
	req.Header.Set("os-version", "")
	req.Header.Set("referer", "https://fanyi.caiyunapp.com/")
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Google Chrome";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	req.Header.Set("x-authorization", "token:qgemv4jr1y38jyq6vhvi")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("bad StatusCode", resp.StatusCode, "body", string(bodyText))
	}
	var dictResponse CaiyunDictResponse
	err = json.Unmarshal(bodyText, &dictResponse)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("--------------------彩云翻译-----------------------")
	fmt.Println(word, "UK:", dictResponse.Dictionary.Prons.En, "US:", dictResponse.Dictionary.Prons.EnUs)
	for _, item := range dictResponse.Dictionary.Explanations {
		fmt.Println(item)
	}
	wg.Done()
}

type BaiDuDictRequest struct {
	Query  string `json:"query"`
	UserID string `json:"user_id"`
}

type BaiDuDictResponse struct {
	TransResult struct {
		Data []struct {
			Dst        string          `json:"dst"`
			PrefixWrap int             `json:"prefixWrap"`
			Result     [][]interface{} `json:"result"`
			Src        string          `json:"src"`
		} `json:"data"`
		From     string `json:"from"`
		Status   int    `json:"status"`
		To       string `json:"to"`
		Type     int    `json:"type"`
		Phonetic []struct {
			SrcStr string `json:"src_str"`
			TrgStr string `json:"trg_str"`
		} `json:"phonetic"`
	} `json:"trans_result"`
	DictResult struct {
		Edict struct {
			Item []struct {
				TrGroup []struct {
					Tr          []string      `json:"tr"`
					Example     []interface{} `json:"example"`
					SimilarWord []string      `json:"similar_word"`
				} `json:"tr_group"`
				Pos string `json:"pos"`
			} `json:"item"`
			Word string `json:"word"`
		} `json:"edict"`
		Collins struct {
			Entry []struct {
				EntryID string `json:"entry_id"`
				Type    string `json:"type"`
				Value   []struct {
					MeanType []struct {
						InfoType string `json:"info_type"`
						InfoID   string `json:"info_id"`
						Example  []struct {
							ExampleID string `json:"example_id"`
							TtsSize   string `json:"tts_size"`
							Tran      string `json:"tran"`
							Ex        string `json:"ex"`
							TtsMp3    string `json:"tts_mp3"`
						} `json:"example"`
					} `json:"mean_type"`
					Gramarinfo []interface{} `json:"gramarinfo"`
					Tran       string        `json:"tran"`
					Def        string        `json:"def"`
					MeanID     string        `json:"mean_id"`
					Posp       []struct {
						Label string `json:"label"`
					} `json:"posp"`
				} `json:"value"`
			} `json:"entry"`
			WordName      string `json:"word_name"`
			Frequence     string `json:"frequence"`
			WordEmphasize string `json:"word_emphasize"`
			WordID        string `json:"word_id"`
		} `json:"collins"`
		From        string `json:"from"`
		SimpleMeans struct {
			WordName  string   `json:"word_name"`
			From      string   `json:"from"`
			WordMeans []string `json:"word_means"`
			Exchange  struct {
				WordPl  []string `json:"word_pl"`
				WordEst []string `json:"word_est"`
				WordEr  []string `json:"word_er"`
			} `json:"exchange"`
			Tags struct {
				Core  []string `json:"core"`
				Other []string `json:"other"`
			} `json:"tags"`
			Symbols []struct {
				PhEn  string `json:"ph_en"`
				PhAm  string `json:"ph_am"`
				Parts []struct {
					Part  string   `json:"part"`
					Means []string `json:"means"`
				} `json:"parts"`
				PhOther string `json:"ph_other"`
			} `json:"symbols"`
		} `json:"simple_means"`
		Lang   string `json:"lang"`
		Oxford struct {
			Entry []struct {
				Tag  string `json:"tag"`
				Name string `json:"name"`
				Data []struct {
					Tag  string `json:"tag"`
					Data []struct {
						Tag   string `json:"tag"`
						P     string `json:"p"`
						PText string `json:"p_text"`
					} `json:"data"`
				} `json:"data"`
			} `json:"entry"`
			Unbox []struct {
				Tag  string `json:"tag"`
				Type string `json:"type"`
				Name string `json:"name"`
				Data []struct {
					Tag     string   `json:"tag"`
					Text    string   `json:"text,omitempty"`
					Words   []string `json:"words,omitempty"`
					Outdent string   `json:"outdent,omitempty"`
					Data    []struct {
						Tag    string `json:"tag"`
						EnText string `json:"enText"`
						ChText string `json:"chText"`
					} `json:"data,omitempty"`
				} `json:"data"`
			} `json:"unbox"`
		} `json:"oxford"`
		Sanyms []struct {
			Tit  string `json:"tit"`
			Type string `json:"type"`
			Data []struct {
				P string   `json:"p"`
				D []string `json:"d"`
			} `json:"data"`
		} `json:"sanyms"`
		Usecase struct {
			Idiom []struct {
				P    string `json:"p"`
				Tag  string `json:"tag"`
				Data []struct {
					Tag  string `json:"tag"`
					Data []struct {
						EnText string `json:"enText"`
						Tag    string `json:"tag"`
						ChText string `json:"chText"`
						Before []struct {
							Tag  string `json:"tag"`
							Data []struct {
								EnText string `json:"enText"`
								Tag    string `json:"tag"`
								ChText string `json:"chText"`
							} `json:"data"`
						} `json:"before,omitempty"`
					} `json:"data"`
				} `json:"data"`
			} `json:"idiom"`
		} `json:"usecase"`
		BaiduPhrase []struct {
			Tit   []string `json:"tit"`
			Trans []string `json:"trans"`
		} `json:"baidu_phrase"`
		QueryExplainVideo struct {
			ID           int    `json:"id"`
			UserID       string `json:"user_id"`
			UserName     string `json:"user_name"`
			UserPic      string `json:"user_pic"`
			Query        string `json:"query"`
			Direction    string `json:"direction"`
			Type         string `json:"type"`
			Tag          string `json:"tag"`
			Detail       string `json:"detail"`
			Status       string `json:"status"`
			SearchType   string `json:"search_type"`
			FeedURL      string `json:"feed_url"`
			Likes        string `json:"likes"`
			Plays        string `json:"plays"`
			CreatedAt    string `json:"created_at"`
			UpdatedAt    string `json:"updated_at"`
			DuplicateID  string `json:"duplicate_id"`
			RejectReason string `json:"reject_reason"`
			CoverURL     string `json:"coverUrl"`
			VideoURL     string `json:"videoUrl"`
			ThumbURL     string `json:"thumbUrl"`
			VideoTime    string `json:"videoTime"`
			VideoType    string `json:"videoType"`
		} `json:"queryExplainVideo"`
	} `json:"dict_result"`
	LijuResult struct {
		Double string   `json:"double"`
		Tag    []string `json:"tag"`
		Single string   `json:"single"`
	} `json:"liju_result"`
	Logid int64 `json:"logid"`
}

func QueryBaiDu(word string, wg *sync.WaitGroup) {
	client := &http.Client{} //定义客户
	// var data = strings.NewReader(`query=good`) //把输入的字符串 data 转换成流 req
	request := BaiDuDictRequest{Query: word}
	buf, err := json.Marshal(request) // Marshal() 返回 request 的 JSON 编码，序列化为一个 byte 数组。
	if err != nil {
		log.Fatal(err) //打印日志，退出程序
	}
	var data = bytes.NewReader(buf) //把输入的 byte 数组转换成流 req
	/* 创建请求 */
	req, err := http.NewRequest("POST", "https://fanyi.baidu.com/langdetect", data)
	if err != nil {
		log.Fatal(err) //打印日志，退出程序
	}
	/* 设置请求头 */
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", `BIDUPSID=2860D678CCB10886990D4D819CC24BAB; PSTM=1603262095; REALTIME_TRANS_SWITCH=1; FANYI_WORD_SWITCH=1; HISTORY_SWITCH=1; SOUND_SPD_SWITCH=1; SOUND_PREFER_SWITCH=1; __yjs_duid=1_3242d9ae5ec151bc4e7e9d908fe0e54c1620968286434; BDUSS=g3Tm93dWtraHZkRUdsdzlJeDEyalp2RmpZOW1RNmlDMUV3d2N0M0c4Z0pqNXhoRVFBQUFBJCQAAAAAAAAAAAEAAABffmtXxKnTsNChutpYSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkCdWEJAnVhfk; BDUSS_BFESS=g3Tm93dWtraHZkRUdsdzlJeDEyalp2RmpZOW1RNmlDMUV3d2N0M0c4Z0pqNXhoRVFBQUFBJCQAAAAAAAAAAAEAAABffmtXxKnTsNChutpYSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkCdWEJAnVhfk; APPGUIDE_10_0_2=1; BAIDUID=3A21E0D3A023E0BF48D5E6AB59882FC3:SL=0:NR=20:FG=1; MAWEBCUID=web_ZbMGPxbtNTJgOhwuypnUUKjyVeChjrWNsCSTdnAIkJyWEtWEYb; BDORZ=FFFB88E999055A3F8A630C64834BD6D0; Hm_lvt_64ecd82404c51e03dc91cb9e8c025574=1652153403,1652175540,1652186130,1652278875; Hm_lpvt_64ecd82404c51e03dc91cb9e8c025574=1652506784; delPer=0; PSINO=1; BAIDUID_BFESS=3A21E0D3A023E0BF48D5E6AB59882FC3:SL=0:NR=20:FG=1; RT="z=1&dm=baidu.com&si=hqolag82j09&ss=l35hyf6s&sl=2&tt=56c&bcn=https%3A%2F%2Ffclog.baidu.com%2Flog%2Fweirwood%3Ftype%3Dperf&ld=5ne&ul=2axq&hd=2b02"; BA_HECTOR=24052g04012lag8kt41h7ulf10q; BDRCVFR[OEHfjv-pq1f]=mk3SLVN4HKm; BDRCVFR[dG2JNJb_ajR]=mk3SLVN4HKm; BDRCVFR[-pGxjrCMryR]=mk3SLVN4HKm; H_PS_PSSID=`)
	req.Header.Set("Origin", "https://fanyi.baidu.com")
	req.Header.Set("Referer", "https://fanyi.baidu.com/?aldtype=16047")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.44")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua", `" Not A;Brand";v="99", "Chromium";v="100", "Microsoft Edge";v="100"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	/* 发起请求 */
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err) //打印日志，退出程序
	}
	defer resp.Body.Close() //defer 会在函数结束后从后往前触发，Close() 手动关闭 Body流，防止内存资源泄露
	/*读取响应*/
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err) //打印日志，退出程序
	}
	if resp.StatusCode != 200 { // 防御式编程，判断状态码是否正确
		log.Fatal("bad StatusCode:", resp.StatusCode, "body", string(bodyText)) //打印日志，退出程序
	}
	// fmt.Printf("%s\n", bodyText)
	var dictResponse BaiDuDictResponse
	err = json.Unmarshal(bodyText, &dictResponse) //Unmarshal()解析 bodyText的 JSON 编码的数据并反序列化将结果存储在 dictResponse 指向的值中。
	if err != nil {
		log.Fatal(err) //打印日志，退出程序
	}
	// fmt.Printf("%#v\n", dictResponse)
	fmt.Println("--------------------百度翻译-----------------------")
	fmt.Println(word, "UK:", dictResponse.DictResult.SimpleMeans.Symbols, "US:", dictResponse.DictResult.SimpleMeans.Symbols)
	for _, item := range dictResponse.DictResult.SimpleMeans.WordMeans {
		fmt.Println(item)
	}

	wg.Done()
}

func main() {
	start := time.Now()
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD example : simpleDict hello`)
		os.Exit(1)
	}
	word := os.Args[1]

	wg := sync.WaitGroup{}
	wg.Add(2)
	go queryCaiyun(word, &wg)
	go QueryBaiDu(word, &wg)

	wg.Wait()
	finish := time.Now()
	fmt.Println("runtime:", finish.Sub(start))
}
