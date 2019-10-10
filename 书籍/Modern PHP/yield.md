#yield

---

		function gen() {
		    yield 'yield1';
		    yield 'yield2';
		}
		 
		$gen= gen();

		var_dump($gen->current());//输出 yield1,send 会额外调用一次 next
		var_dump($gen->send('ret1')); //输出是 yield2 