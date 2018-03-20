<div class="panel panel-default">
  <div class="panel-heading">
    作者信息
  </div>
  <div class="panel-body">
    <div class="media">
      <div class="media-left">
        <a href="/user/{{.Topic.User.Username}}">
          <img src="{{.Topic.User.Image}}" title="{{.Topic.User.Username}}" class="avatar">
        </a>
      </div>
      <div class="media-body">
        <div class="media-heading">
          <a href="/user/{{.Topic.User.Username}}">{{.Topic.User.Username}}</a>
        </div>
        {{if .Topic.User.URL}}<a href="{{.Topic.User.URL}}" target="_blank">{{.Topic.User.URL}}</a>{{end}}
      </div>
      {{if .Topic.User.Signature}}
      <div style="color: #7A7A7A; font-size: 12px; margin-top:5px;">
        <i>“ {{.Topic.User.Signature}} ” </i>
      </div>
      {{end}}
    </div>
  </div>
</div>