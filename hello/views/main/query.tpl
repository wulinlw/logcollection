<h5>查询数据源：{{.tableName}}</h5>
<h5>查询时间：{{.params.Stime}} -- {{.params.Etime}}</h5>
{{if .list}}
{{with $value := index .list 0}}
   <h5>数据源路径：{{printf .File_name}}</h5>
{{end}}
<h5>共{{.totalItem}}记录</h5>
{{/*index .list 0 */}}
<div class="default">
    <div>{{.pagerhtml}}</div>
</div>
<div class="col-md-12">
    <table class="table">
      <caption></caption>
      <thead>
        <tr>
          <th>id</th>
		  <th>time</th>
		  <th>line</th>
		  <th>content</th>
        </tr>
      </thead>
      <tbody>
		{{range $key, $val := .list}}
        <tr>
          <th>{{$val.Id}}</th>
		  <td>{{$val.Crtime | showTime}}</td>
		  <td>{{$val.Line}}</td>
		  <td>{{$val.Content}}</td>
        </tr>
		{{end}} 
      </tbody>
    </table>
</div>
{{else}}
没有找到记录
{{end}}






