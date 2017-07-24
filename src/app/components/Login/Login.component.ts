import { Component, OnInit } from '@angular/core';
import { AccountAPI } from '../../services/AccountAPI.service';

@Component({
  selector: 'login',
  templateUrl: './Login.component.html',
  styleUrls: ['./Login.component.css']
})
export class LoginComponent implements OnInit { 
  constructor(private accountAPI: AccountAPI) {}

  login() {
    //TODO: Send token requirements from the form
    this.accountAPI.Token('','').subscribe();
    localStorage.setItem('dwn_token','');
  }

  ngOnInit() {
  }
}