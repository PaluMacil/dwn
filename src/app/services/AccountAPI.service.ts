import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class AccountAPI {
  constructor(private http: HttpClient) {}

  session?: Session = null;

  public Token(email: string, password: string): Observable<Session> {
    return this.http.post<Session>('/api/account/token', { email, password })
  }
}

export interface Session {
			Token:     string;
			User:      User;
			CreatedAt: string;
			HeartBeat: string;
    }
    
export interface User {
	ID:        number;
	Role:      number;
	Email:     string;
	Name:      string;
	CreatedAt: string;
}

export interface UserInfo {
  ID:      number;
  Name:    string;
  IsAdmin: boolean;
}