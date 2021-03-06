<div class="row">
  <div class="col-md-9">
    <div class="panel panel-default">
      <div class="panel-heading">
        <a href="/">主页</a> / 编辑话题
      </div>
      <div class="panel-body">
        {{template "../components/flash_error.tpl" .}}
        <form method="post" action="/topic/edit/{{.Topic.Id}}" enctype="multipart/form-data">
          <div class="form-group">
            <label for="title">标题</label>
            <input type="text" class="form-control" Id="title" name="title" value="{{.Topic.Title}}" placeholder="标题">
          </div>
          <div class="form-group">
            <label for="title">内容</label>
            <div id="my-editormd" >
              <textarea id="my-editormd-markdown-doc" name="my-editormd-markdown-doc" style="display:none;">{{.Topic.Content}}</textarea>
              <!-- 注意：name属性的值-->
              <textarea id="my-editormd-html-code" name="my-editormd-html-code" style="display:none;"></textarea>
            </div>
          </div>
          <div class="form-group">
            <label for="title">标签</label>
             <div>
              {{range .Tags}}
                <input type="checkbox" name="tids" value="{{.Id}}" id="tag_{{.Id}}">
                <label for="{{.Id}}">{{.Name}}</label>&nbsp;
              {{end}}
            </div>
          </div>
           <div class="form-group">
            <label >选择文件</label>
            <input type="file" class="form-control" name="file">
          </div>
          <button type="submit" class="btn btn-default">发布</button>
        </form>
      </div>
    </div>
  </div>
</div>
<script type="text/javascript">
  $(function () {
    {{range .TopicTags}}
    $("#tag_{{.Id}}").attr("checked", true);
    {{end}}
  });
</script>

<script type="text/javascript">
  $(function() {
      editormd("my-editormd", {//注意1：这里的就是上面的DIV的id属性值
          width   : "90%",
          height  : 450,
          syncScrolling : "single",
          path    : "/static/editor/lib/",//注意2：你的路径
          saveHTMLToTextarea : true//注意3：这个配置，方便post提交表单
      });
  });
</script>