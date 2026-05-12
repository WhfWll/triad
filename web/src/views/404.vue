<template>
	<div class="content">
		<div>
			<div class="errorbg">
				<div>
					<strong>{{code}}</strong>
					<div>抱歉，您访问的页面不存在</div>
				</div> 
			</div>
			<el-button type="primary" style="margin-top:30px;"  @click="returnback()" >返回首页</el-button>
		</div>
	</div>
</template>
<style>
	.content{
		background: #f5f6fa;
	}
	.content >div{
		position: absolute;
		top: 50%;
		left: 50%;
		margin-top: -300px;
		margin-left: -146px;
		text-align: center;
	}
	.errorbg{ 
		width: 292px;
		height: 241px; 
		background:url(../assets/Error/error.png) no-repeat; 
	}
	.errorbg >div{
		width: 256px;
		height: 128px;
		margin-left: 15px;
		padding-top: 100px; 
	}
	.errorbg >div *{
		display: block;
		text-align: center;
	}
	.errorbg >div strong{
		font-size: 52px;
		padding-top:20px;
		margin-bottom: 10px;
	}
	.errorbg >div  div{
		font-size: 12px;
		color: #6c706f;
	}
</style>
<script>

export default {
    data () { 
        return { 
			code:'404',
			role:0, 
        } 
	}, 
	created() {
		this.role = this.commonjs.decryptCBC(localStorage.getItem('role'), this.commonjs.myKey); 
	},
    methods: { 
    	returnback:function(){
			if(this.role == 3){
				this.$router.push('/log');
			}else{
				this.$router.push('/');
			}
            
        },   
    }
}
</script>
