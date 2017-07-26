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
  loginErr?: string = null;

  login() {
    this.accountAPI.Token(this.email,this.password).subscribe(
      data => {
        localStorage.setItem('dwn_token',data.Token);
        this.accountAPI.session = data;
        $('#login-modal').modal('hide');
        this.email = '';
        this.password = '';
        this.loginErr = null;
      },
      err => {
        this.loginErr = err.message;
        this.password = '';
      }
    );
  }

  logout() {
    localStorage.removeItem('dwn_token');
    this.accountAPI.session = null;
    this.accountAPI.Logout().subscribe(
      () => { console.log('Logged out') },
      (err) => { console.log('Could not log out. ' + err.message) }
    );
  }

  ngOnInit() {
  }
}