.text
	in r16,0x04 ; ieliek I/O porta 0x04 datus reģistrā R16
	ori r16,0x20 ; ???? ???? OR 0x20, rezultāts iet iekš R16
	out 0x04,r16 ; ieliek reģistra R16 datus I/O portā 0x04

	loop:	
		in r16,0x05 ;
		andi r16,0xDF ; ???? ???? AND 0xDF, rezultāts iet iekš R16
		out 0x05,r16 ;

		; 16000000/196605 ~= 81
		ldi r16,78 ; ieliek 78 iekš R16
		delay1:
			; delay_inner1 ~= 196605c
			ldi r24,0xFF ; (1c)
			ldi r25,0xFF ; (1c)
			delay_inner1:
				; atņem 1 no reģistru pāra r24 un r25, rezultāts iet iekš pāra
				; ja rezultāts ir 0, statusa reģistrā (SREG) iestata Z (zero) bitu uz HIGH
				sbiw r24,1 ; (2c)
				brne delay_inner1 ; ja Z nav HIGH, lec uz delay_inner1 () (1c ja false, 2c ja true)
			subi r16,1 ; atņem 1, rezultāts iekš R16 (1c)
			brne delay1 ; (1/2c)

		in r16,0x05
		ori r16,0x20
		out 0x05,r16

		ldi r16,78
		delay2:
			ldi r24,0xFF
			ldi r25,0xFF
			delay_inner2:
				sbiw r24,1
				brne delay_inner2
			subi r16,1
			brne delay2

		jmp loop
