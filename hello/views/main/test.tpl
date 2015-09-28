
<div class="row col-sm-10">
<form class="form-horizontal" action="/main/query">
  <div class="form-group">
    <label for="logId" class="col-sm-1 control-label">日志</label>
    <div class="col-sm-10">
      <select class="form-control" name="id">
			{{range $key, $val := .apps}}
			<option value="{{$val.id}}">{{$val.describe}}</option>
			{{end}} 
	  </select>
    </div>
  </div>

  <div class="form-group">
	    <label for="stime" class="col-sm-1 control-label">时间</label>
	  <div class="col-lg-3">
	    <div class="input-group">
	      <input type="text" class="form-control" onfocus="WdatePicker({dateFmt:'yyyy-MM-dd HH:mm:ss'})" name="stime">
	    </div><!-- /input-group -->
	  </div><!-- /.col-lg-2 -->
	  <div class="col-lg-3">
	    <div class="input-group">
	      <input type="text" class="form-control" onfocus="WdatePicker({dateFmt:'yyyy-MM-dd HH:mm:ss'})" name="etime">
	    </div><!-- /input-group -->
	  </div><!-- /.col-lg-2 -->
  </div>

  <div class="form-group">
    <div class="col-sm-offset-1 col-sm-10">
      <button type="submit" class="btn btn-default">submit</button>
    </div>
  </div>
</form>
</div>





