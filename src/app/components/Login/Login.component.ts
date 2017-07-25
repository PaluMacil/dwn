import { Component, OnInit } from '@angular/core';
import { AccountAPI } from '../../services/AccountAPI.service';

@Component({
  selector: 'login',
  templateUrl: './Login.component.html',
  styleUrls: ['./Login.component.css']
})
export class LoginComponent implements OnInit { 
  constructor(private accountAPI: AccountAPI) {}

  email: string;
  password: string;

  login() {
    this.accountAPI.Token(this.email,this.password).subscribe(
      data => {
        localStorage.setItem('dwn_token',data.Token);
        this.accountAPI.session = data;
      },
      err => {
        console.log('Could not log in: ' + JSON.stringify(err));
      }
    );
    
  }

  ngOnInit() {
  }
}