package logenc

/*func TestScan_procFileSearch(t *testing.T) {
	wg := sync.WaitGroup{}
	ch := make(chan string, 3)
	var mutex sync.Mutex
	countL := 0

	scan := Scan{}
	scan.LimitResLines = 5
	scan.ChRes = make(chan string, 3)
	scan.Text = "aaa"

	go func() { // reader
	ext:
		for {
			select {
			case line, ok := <-scan.ChRes:
				if !ok {
					break ext
				}
				fmt.Println(line)
			}
		}
	}()

	for i := 0; i < 2; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			//End:
		ext:
			for {
				select {
				case line, ok := <-ch:
					//mutex.Lock()
					if !ok {
						break ext
					}

					if !ok && scan.LimitResLines == 0 {
						break ext
					} else if scan.LimitResLines <= countL && scan.LimitResLines != 0 {
						//cancel()
						//close(scan.ChRes)
						break
					}

					csvSt, str := line, line
					if strings.Contains(str, scan.Text) {
						counterES.Add(1)
						scan.ChRes <- csvSt
						mutex.Lock()
						countL++
						mutex.Unlock()

					}
					counter.Add(1)
					//return
				}
			}
			//close(ch)
		}()
	}

	for i := 0; i < 20; i++ {
		ch <- string("aaa")
	}

	// if countL >= scan.LimitResLines || scan.LimitResLines == 0 {
	// 	close(ch)
	// }

	//else if t.LimitResLines == 0 {
	close(ch)
	//}

	wg.Wait()

}//*/
