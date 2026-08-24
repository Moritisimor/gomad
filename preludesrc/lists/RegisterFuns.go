package lists

import (
	"github.com/Moritisimor/gomad/eval"
	"github.com/Moritisimor/gomad/value"
)

func RegisterFuns(env *value.Env) error {
	if _, err := eval.DoString(`
		(letfun foldl (f acc l)
    		(do
      			(letfun aux (a h t)
        			(if (= unit t)
          			a
          			(aux (f a h) (car t) (cdr t))))
      
    			(aux acc (car l) (cdr l))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun begins_with (l1 l2)
    		(if (< (len l1) (len l2))
      			false
      			(do
        			(letfun aux (l1h l1t l2h l2t)
          			(if (= unit l2t)
            			true
            			(if (= l1h l2h)
              				(aux (car l1t) (cdr l1t) (car l2t) (cdr l2t))
              				false)))
                
        			(aux (car l1) (cdr l1) (car l2) (cdr l2)))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun ends_with (l1 l2) (begins_with (rev l1) (rev l2)))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun list_init (n f)
    		(do
      			(letfun aux (acc i)
        			(if (< i 0)
          			acc
          			(aux (cons (f i) acc) (dec i))))
          
      			(aux () (dec n))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun rev (l)
    		(do
      			(letfun aux (acc h t)
        			(if (= unit t)
          			acc
          			(aux (cons h acc) (car t) (cdr t))))
          
      			(aux () (car l) (cdr l))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun map (f l)
    		(do
      			(letfun aux (acc h t)
        			(if (= unit t)
          				(rev acc)
          				(aux (cons (f h) acc) (car t) (cdr t))))
        
      			(aux () (car l) (cdr l))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun mapi (f l)
    		(do
      			(letfun aux (acc h t i)
        			(if (= unit t)
          				(rev acc)
          				(aux (cons (f h i) acc) (car t) (cdr t) (inc i))))
          
      			(aux () (car l) (cdr l) 0)))
	`, env); err != nil { 
		return err 
	}

	if _, err := eval.DoString(`
		(letfun filter (f l)
    		(do
      			(letfun aux (acc h t)
        			(if (= unit t)
          				(rev acc)
          				(if (f h)
            				(aux (cons h acc) (car t) (cdr t))
            				(aux acc (car t) (cdr t)))))
        
      			(aux () (car l) (cdr l))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun foreach (f l)
    		(do
      			(letfun aux (h t)
        			(do
          				(if (= unit t)
          				unit
          				(do
            				(f h)
            				(aux (car t) (cdr t))))))
          
      			(aux (car l) (cdr l))))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun foreachi (f l)
    		(do
      			(letfun aux (h t i)
        			(do
          				(if (= unit t)
            			unit
            			(do 
              				(f h i)
              				(aux (car t) (cdr t) (inc i))))))
          
      			(aux (car l) (cdr l) 0)))
	`, env); err != nil {
		return err
	}

	if _, err := eval.DoString(`
		(letfun range (start end list)
    		(do
      			(letfun aux (acc h t i)
        			(if (isunit t)
          				(rev acc)
          				(if (and (>= i start) (<= i end))
            				(aux (cons h acc) (car t) (cdr t) (inc i))
            				(aux acc (car t) (cdr t) (inc i)))))
        
      			(aux () (car list) (cdr list) 0)))
	`, env); err != nil {
		return err
	}

	return nil
}
