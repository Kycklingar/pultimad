package yp

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

//const postHtml = `
//	<div class="col s12 m6">
//		<div class="card large yp-post" id="p30255357">
//			<div class="card-image waves-effect waves-block waves-light">
//				<img class="activator lazyload" src="/static/img/file_placeholder.png" data-src="/patreon_data/2533136/30255357/thumb.jpg" alt="Some doodles">
//			</div>
//			<div class="card-content">
//				<span class="card-title activator grey-text text-darken-4">Some doodles<i class="material-icons right">more_vert</i></span>
//				<br><span style="font-size: 85%;" class="grey-text post-time">2019-09-26T09:17:51+00:00</span>
//			</div>
//			<div class="card-action">
//				<a href="/patreon_data/2533136/30255357/koraremonmegamimon.png" target="_blank">Post file</a>
//				<i class="material-icons right yp-flag" title="Flag this post" data-pid="30255357">flag</i>
//			</div>
//			<div class="card-reveal">
//				<span class="card-title grey-text text-darken-4">Some doodles <small style="font-size: 60%;" class="post-time">2019-09-26T09:17:51+00:00</small><i class="material-icons right">close</i></span>
//				<div class="post-body"><p>Some doodling I got around to doing yesterday. One of Akahimon holding a certain key and bit of Kõraremonmegamimon (name pending?) concept and showing off a bit.</p></div>
//				<div class="card-attachments">
//					<hr style="height:1px;border:none;color:#212121;background-color:#212121;">
//					<span class="card-title grey-text text-darken-4">Media (3)</span>
//					<p style="margin-top:-15px;">
//					<br><a href="https://yiff.party/patreon_data/2533136/30255357/42810327/koraremonmegamimon.png" target="_blank">koraremonmegamimon.png</a> (1.6MiB)
//					<br><a href="https://yiff.party/patreon_data/2533136/30255357/42810345/koraremon_solarcannon.png" target="_blank">koraremon_solarcannon.png</a> (910.8KiB)
//					<br><a href="https://yiff.party/patreon_data/2533136/30255357/42810350/akahimon_key.png" target="_blank">akahimon_key.png</a> (1017.5KiB)
//					</p>
//				</div>
//			</div>
//		</div>
//	</div>
//`

const postHtml = `
<div class="col s12 m6">
<div class="card large yp-post" id="p30825495">
<div class="card-image waves-effect waves-block waves-light">
<img class="activator lazyload" src="/static/img/file_placeholder.png" data-src="https://data.yiff.party/patreon_data/2533136/30825495/thumb.jpg" alt="Inktober 14-15">
</div>
<div class="card-content">
<span class="card-title activator grey-text text-darken-4">Inktober 14-15<i class="material-icons right">more_vert</i></span>
<br><span style="font-size: 85%;" class="grey-text post-time">2019-10-17T09:54:47+00:00</span>
</div>
<div class="card-action">
<a href="https://data.yiff.party/patreon_data/2533136/30825495/inktober_14.png" target="_blank">Post file</a>
<i class="material-icons right yp-flag" title="Flag this post" data-pid="30825495">flag</i>
</div>
<div class="card-reveal">
<span class="card-title grey-text text-darken-4">Inktober 14-15 <small style="font-size: 60%;" class="post-time">2019-10-17T09:54:47+00:00</small><i class="material-icons right">close</i></span>
<div class="post-body"><p>Inktobers 14-15. A little late due my night shifts and art classes but here we go.</p><p>day 14: Street Walking - Idena asks the tribal 'punny what she's doing. The reply makes her want to know if she's compatible too...<br/>day 15: First Time - It is tribal 'pun's first time, pretty much in what it resulted as well.<br/></p></div>
<div class="card-attachments">
<hr style="height:1px;border:none;color:#212121;background-color:#212121;">
<span class="card-title grey-text text-darken-4">Media (2)</span>
<p style="margin-top:-15px;">
<br><a href="https://data.yiff.party/patreon_data/2533136/30825495/44491337/inktober_14.png" target="_blank">inktober_14.png</a> (376.2KiB)
<br><a href="https://data.yiff.party/patreon_data/2533136/30825495/44491347/inktober_15.png" target="_blank">inktober_15.png</a> (295.8KiB)
</p>
</div>
<div class="card-comments">
<span class="card-title grey-text text-darken-4 yp-post-comment-title">Comments (2)</span>
<div class="yp-post-comment">
<div>
<img class="yp-post-comment-avatar lazyload" data-src="https://data.yiff.party/avatars/patreon/thumb/12461443.jpg" src="/static/img/avatar_placeholder.png" alt="user avatar">
</div>
<div><span class="yp-post-comment-head"><a href="https://www.patreon.com/user?u=12461443" rel="noreferrer" target="_blank"><strong>User #12461443</strong></a> - <span class="yp-post-comment-time" data-utc="1571325475">17 Oct 19 15:17</span></span>
<div class="yp-post-comment-body">Bunny buries a bone, what happens in future will surprise you!</div></div>
</div>
<div class="yp-post-comment">
<div>
<img class="yp-post-comment-avatar lazyload" data-src="https://data.yiff.party/avatars/patreon/thumb/325476.jpg" src="/static/img/avatar_placeholder.png" alt="user avatar">
</div>
<div><span class="yp-post-comment-head"><a href="https://www.patreon.com/user?u=325476" rel="noreferrer" target="_blank"><strong>User #325476</strong></a> - <span class="yp-post-comment-time" data-utc="1571353160">17 Oct 19 22:59</span></span>
<div class="yp-post-comment-body">https://youtu.be/S6Xdfql7C0Y</div></div>
</div>
</div>
</div>
</div>
</div>
`

func TestPostParse(t *testing.T) {
	n, err := html.Parse(strings.NewReader(postHtml))
	var r func(no *html.Node) bool
	r = func(no *html.Node) bool {
		if no.Type == html.ElementNode && no.Data == "div" {
			n = no
			return true
		}
		for c := no.FirstChild; c != nil; c = c.NextSibling {
			if r(c) {
				return true
			}
		}
		return false
	}
	r(n)

	if err != nil {
		t.Fatal(err)
	}
	post, err := parsePost(&node{*n})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(post.ID)
	fmt.Println(post.Title)
	fmt.Println(post.Body)
	fmt.Println(post.FileURL)
	fmt.Println(post.Attachments)
}
