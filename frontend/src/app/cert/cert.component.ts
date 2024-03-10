import { Component, OnInit } from '@angular/core';
import { CertificateService } from './cert-service';
import { KeycloakService } from 'keycloak-angular';
import { CertificateType } from './cert-data';
import { saveAs } from "file-saver";

@Component({
  selector: 'app-cert',
  templateUrl: './cert.component.html',
  styleUrl: './cert.component.css'
})
export class CertComponent implements OnInit{
  certs: CertificateType[] = [];
  loading: boolean = false;
  constructor(private cs: CertificateService, 
              private kc: KeycloakService){}
  ngOnInit(): void {
    this.readCerts();
  }
  readCerts(): void{
    this.loading = true;
    this.cs.getCerts().subscribe(data => {
      this.certs = data;
      this.loading = false;
    });
  }
  create():void{
    this.cs.createCert().subscribe(() => this.readCerts());
    
  }
  download(cn: string):void{
    this.cs.getCertByCN(cn).subscribe(data => saveAs(data, cn + ".ovpn"));    
  }
  revoke(cn:string):void{
    this.cs.revokeByCN(cn).subscribe(()=>this.readCerts());
  }
  
}
