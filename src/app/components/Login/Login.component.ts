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

  ngOnInit() {
    let login = this;
    $("#login-modal").on("hidden.bs.modal", function () {
      login.email = '';
      login.password = '';
      login.loginErr = null;
    });
  }
}